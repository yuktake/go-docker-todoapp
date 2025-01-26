FROM golang:1.23.5

WORKDIR /usr/local/go/src/go_crm/

RUN go mod init github.com/yuktake/go_crm

# airをインストールし、コンテナ起動時にホットリロードできるよう実行する
# こっちのライブラリだと動かない（github.com/cosmtrek/air）
RUN go install github.com/air-verse/air@latest

COPY application/. /usr/local/go/src/go_crm/

# ローカルでgolang:1.21.1で動かしていたものをgo.modのバージョンを変更して、go mod tidyを実行している
RUN go mod tidy
RUN go mod download