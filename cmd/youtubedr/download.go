package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"strconv"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Downloads a video from youtube",
	Example: `youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(download(args[0]))
	},
}

var (
	ffmpegCheck error
	outputFile  string
	outputDir   string
	itagString  string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputFile, "filename", "o", "", "The output file, the default is genated by the video title.")
	downloadCmd.Flags().StringVarP(&outputDir, "directory", "d", ".", "The output directory.")
	downloadCmd.Flags().StringVarP(&itagString, "itag", "i", "", "Itag number of the stream.")
	addQualityFlag(downloadCmd.Flags())
	addMimeTypeFlag(downloadCmd.Flags())
}

func download(id string) error {
	video, format, err := getVideoWithFormat(id)
	if err != nil {
		return err
	}

	log.Println("download to directory", outputDir)

        if itagString != "" {
                itagNo, err := strconv.Atoi(itagString)
                if err != nil {
                        return fmt.Errorf("Invalid Itag number %s", itagString)
                }
                return downloader.DownloadByItag(context.Background(), outputFile, video, itagNo)
        }

	if strings.HasPrefix(outputQuality, "hd") {
		if err := checkFFMPEG(); err != nil {
			return err
		}
		return downloader.DownloadComposite(context.Background(), outputFile, video, outputQuality, mimetype)
	}

	return downloader.Download(context.Background(), video, format, outputFile)
}

func checkFFMPEG() error {
	fmt.Println("check ffmpeg is installed....")
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		ffmpegCheck = fmt.Errorf("please check ffmpegCheck is installed correctly")
	}

	return ffmpegCheck
}
