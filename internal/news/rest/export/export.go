package export

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Benzinga/sdk-go/benzinga"
	"github.com/Benzinga/sdk-go/internal/cli"
	"github.com/Benzinga/sdk-go/internal/cli/prompt"
	"github.com/Benzinga/sdk-go/pkg/client/rest/news"
)

const (
	DateFormat = "2006_01_02"
	PageSize   = 100
)

type Config struct {
	Global             *cli.Config
	StartDate, EndDate time.Time
	UseGzip            bool
	OutputDirectory    string
}

func NewConfig(global *cli.Config) *Config {
	return &Config{
		Global: global,
	}
}

func DatePath(prefix string, d time.Time) string {
	prefix = DateDirectory(prefix, d)

	filename := d.Format(DateFormat) + ".json.gz"

	return path.Join(prefix, filename)
}

func DateDirectory(prefix string, d time.Time) string {
	year := strconv.Itoa(d.Year())
	month := d.Month().String()

	return path.Join(prefix, year, month)
}

func handleDate(ctx context.Context, client *benzinga.Client, dirPrefix, token string, date time.Time) error {
	directory := DateDirectory(dirPrefix, date)

	CreateDirIfNotExist(directory)

	request := client.News()
	request.SetAPIToken(token)
	request.SetPageSize(PageSize)
	request.SetDate(date)
	request.SetSortDirection(news.Ascending)
	request.SetSortField(news.CreatedField)
	request.SetDisplayOutput(news.FullOutput)

	filename := DatePath(dirPrefix, date)

	// check if file exists for date and skip if so
	if _, err := os.Stat(filename); err == nil {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		zap.L().Fatal("create file error", zap.Error(err), zap.String("filename", filename))
	}

	gw := gzip.NewWriter(file)

	bw := bufio.NewWriter(gw)

	writer := json.NewEncoder(bw)

	var fileWrittenTo bool

	defer func() {
		bw.Flush()
		gw.Close()

		file.Close()

		if !fileWrittenTo {
			os.Remove(filename)
		}
	}()

	for page := 0; page < 1000; page++ {
		request.SetPage(page)

		u, err := request.URL()
		if err != nil {
			zap.L().Error("request url error", zap.Error(err))

			return err
		}

		results, err := request.Exec(ctx)
		if err != nil {
			zap.L().Error("request error", zap.Error(err), zap.Stringer("url", u))

			return err
		}

		for _, result := range results {
			if err := writer.Encode(&result); err != nil {
				zap.L().Error("write json error", zap.Error(err), zap.Int("nid", result.ID))

				return err
			}

			fileWrittenTo = true
		}

		zap.L().Info("retrieved", zap.String("date", date.Format(news.DateFormat)), zap.Int("page", page), zap.Int("results", len(results)))

		if len(results) < PageSize {
			break
		}
	}

	return nil
}

func Start(config *Config) {
	client := benzinga.NewClient(nil)

	token, err := prompt.ForStringMasked("API Token")
	if err != nil {
		zap.L().Error("API token error", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	exportDirectoryPrefix := getDirectory()

	handleYear := func(year int) error {
		date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

		endOfYear := date.AddDate(1, 0, 0)

		zap.L().Sugar().Infof("Starting with %s to directory prefix %s", date.Format(news.DateFormat), exportDirectoryPrefix)

		for !endOfYear.Equal(date) && !date.After(time.Now()) {
			if err := handleDate(ctx, client, exportDirectoryPrefix, token, date); err != nil {
				zap.L().Error("handle date error", zap.String("date", date.Format(news.DateFormat)), zap.Error(err))

				return err
			}

			date = date.AddDate(0, 0, 1)
		}

		return nil
	}

	var wg sync.WaitGroup

	for year := 2011; year < time.Now().Year()+1; year++ {
		wg.Add(1)

		go func(y int) {
			if err := handleYear(y); err != nil {
				zap.L().With(zap.Error(err)).Sugar().Fatalf("handle year %d error", y)
			}
			wg.Done()
		}(year)
	}

	wg.Wait()

	cancel()

	zap.L().Info("exiting.")
}

func getDirectory() string {
	val, err := prompt.ForString("Export Directory")
	if err != nil {
		zap.L().Fatal("output directory error", zap.Error(err))

		return ""
	}

	return val
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		zap.L().Sugar().Infof("creating directory %s", dir)

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			zap.L().Error("error creating directory", zap.Error(err), zap.String("dir", dir))
			os.Exit(1)
		}
	}
}
