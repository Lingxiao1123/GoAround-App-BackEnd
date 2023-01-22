package backend

import (
	"context"
	"fmt"   //打印输出
	"io"    //上传文件

	"around/constants"  // gcs bucket

	"cloud.google.com/go/storage"
)

var (
    GCSBackend *GoogleCloudStorageBackend
)


//storage Client
//与electic search不同，需要知道folder的名字（Bucket）
type GoogleCloudStorageBackend struct {
    client *storage.Client
    bucket string
}

//receiver
//objectName: === bucket name：知道文件存在哪里
//返回值：string（返回存好的链接）
// func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
//     ctx := context.Background()
//     object := backend.client.Bucket(backend.bucket).Object(objectName) //远程创建bucket
//     wc := object.NewWriter(ctx)    //创建writter
//     if _, err := io.Copy(wc, r); err != nil {
//         return "", err      //如果出现error，返回空连接
//     }

//     if err := wc.Close(); err != nil {
//         return "", err
//     }

// 	//object.ACL() --- > 设置reader权限
//     if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
//         return "", err
//     }

//     attrs, err := object.Attrs(ctx) //Attrs：Attributes
//     if err != nil {
//         return "", err
//     }

//     fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
//     return attrs.MediaLink, nil
// 	//返回url
// }

func InitGCSBackend() {
    client, err := storage.NewClient(context.Background())
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
        client: client,
        bucket: constants.GCS_BUCKET,
    }
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
    ctx := context.Background()
    object := backend.client.Bucket(backend.bucket).Object(objectName)
    wc := object.NewWriter(ctx)
    if _, err := io.Copy(wc, r); err != nil {
        return "", err
    }

    if err := wc.Close(); err != nil {
        return "", err
    }

    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
    }

    attrs, err := object.Attrs(ctx)
    if err != nil {
        return "", err
    }

    fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}


// //struct constucter：创建backend ---> new client 以及获得bucket
// func InitGCSBackend() {
// 	client, err := storage.NewClient(context.Background())  //initiate new client // context:request deadline
// 	if err != nil{
// 		panic(err)
// 	}

// 	GCSBackend = &GoogleCloudStorageBackend{
//         client: client,
//         bucket: constants.GCS_BUCKET,
//     }
// }


