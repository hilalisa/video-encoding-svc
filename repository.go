package main

import (
	"context"
	"database/sql"
	pb "github.com/agxp/cloudflix/video-encoding-svc/proto"
	"github.com/minio/minio-go"
	"github.com/opentracing/opentracing-go"
	"log"
	"github.com/go-redis/redis"
	"github.com/xfrr/goffmpeg/transcoder"
	"strings"
	"strconv"
	"fmt"
)

type Repository interface {
	Encode(p opentracing.SpanContext, video_id string) (*pb.Response, error)
}

type EncodeRepository struct {
	s3     *minio.Client
	pg     *sql.DB
	cache  *redis.Client
	tracer *opentracing.Tracer
}

func (repo *EncodeRepository) Encode(parent opentracing.SpanContext, video_id string) (*pb.Response, error) {
	sp, _ := opentracing.StartSpanFromContext(context.TODO(), "Encode_Repo", opentracing.ChildOf(parent))

	sp.LogKV("video_id", video_id)

	defer sp.Finish()

	// if already encoded, exit
	// else:

	// get filepath
	psSP, _ := opentracing.StartSpanFromContext(context.TODO(), "PG_EncodeGetFilePath", opentracing.ChildOf(sp.Context()))

	var file_path string

	selectQuery := `select file_path from videos where id=$1`
	err := repo.pg.QueryRow(selectQuery, video_id).Scan(&file_path)
	if err != nil {
		log.Print(err)
		psSP.Finish()
		return nil, err
	}
	psSP.Finish()

	sp.LogKV("file_path", file_path)

	video_path := "/tmp/" + video_id

	// pull file from S3
	err = repo.s3.FGetObject("videos", file_path, video_path, minio.GetObjectOptions{})
	if err != nil {
		log.Print(err)
		psSP.Finish()
		return nil, err
	}
	// generate thumbnail
	thumb_path := video_path + ".jpg"

	var resolution string

	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input file path and output file path
	err = trans.Initialize( video_path, thumb_path )
	// Handle error...
	if err != nil {
		log.Println(err)
	}

	trans.MediaFile().SetFrameRate(1)
	trans.MediaFile().SetSeekTimeInput("00:00:04")
	trans.MediaFile().SetDurationInput("00:00:01")
	resolution = trans.MediaFile().Resolution()

	// Start transcoder process
	done, err := trans.Run()
	if err != nil {
		log.Println(err)
	}

	progress, err := trans.Output()
	if err != nil {
		log.Println(err)
	}

	for msg := range progress {
		log.Println(msg)
	}

	// This channel is used to wait for the process to end
	<-done

	//var width float64
	var height float64

	if resolution != "" {
		resolution := strings.Split(resolution, "x")
		if len(resolution) != 0 {
			//width, _ = strconv.ParseFloat(resolution[0], 64)
			height, _ = strconv.ParseFloat(resolution[1], 64)
		}
	}




	if height > 144 {
		// encode 144p
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_144.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}
		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("256x144")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}

	if height > 240 {
		// encode 240p
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_240.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}
		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("426x240")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}

	if height > 360 {
		// encode 360p
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_360.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}
		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("640x360")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}

	if height > 480 {
		// encode 480p
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_480.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}
		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("854x480")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}

	if height > 720 {
		// encode 720
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_720.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}
		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("1280x720")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}

	if height > 1080 {
		// encode 1080p
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( video_path, video_path + "_1080.mp4" )
		// Handle error...
		if err != nil {
			log.Println(err)
		}

		trans.MediaFile().SetVideoCodec("libx264")
		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetResolution("1920x1080")
		trans.MediaFile().SetQuality(26)
		trans.MediaFile().SetFrameRate(30)

		// Start transcoder process
		done, err := trans.Run()
		if err != nil {
			log.Println(err)
		}

		progress, err := trans.Output()
		if err != nil {
			log.Println(err)
		}

		for msg := range progress {
			log.Println(msg)
		}

		// This channel is used to wait for the process to end
		<-done
	}


	return nil, nil
}
