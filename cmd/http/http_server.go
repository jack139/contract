package http

import (
	"fmt"
	"log"
	"os"
	"time"
	"errors"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/Ferluci/fast-realip"
	"github.com/spf13/cobra"
)


var (
	output = log.New(os.Stdout, "", 0)
	//userKeyfilePath string
)

// "github.com/AubSs/fasthttplogger"
func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}
// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url> - <status> - <response-time us> | <user-agent>
// [2017/05/31 - 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello - 200 - 48.279µs | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func combined(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v (%v) | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			realip.FromRequest(ctx),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}


/* 入口 */
func runServer(port string/*, userPath string*/) {
	// 装入用户appid和secret
	//err := loadSecretKey(userPath)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 保存到全局变量，注册新用户时要用
	//userKeyfilePath = userPath

	/* router */
	r := router.New()
	r.GET("/", index)
	r.POST("/api/test", doNonthing)
	//r.POST("/api/query_deals", queryDeals)
	//r.POST("/api/query_by_assets", queryByAsstes)
	//r.POST("/api/query_block", queryBlock)
	//r.POST("/api/query_raw_block", queryRawBlock)
	//r.POST("/api/biz_contract", bizContract)
	//r.POST("/api/biz_delivery", bizDelivery)
	//r.POST("/api/biz_register", bizRegister)


	fmt.Printf("start HTTP server at 0.0.0.0:%s\n", port)

	/* 启动server */
	s := &fasthttp.Server{
		Handler: combined(r.Handler),
		Name: "FastHttpLogger",
	}
	log.Fatal(s.ListenAndServe(":"+port))
}


/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}


func HttpCliCmd() *cobra.Command {
	cmd := &cobra.Command{	// 启动http服务
		Use:   "http <port>",
		Short: "start http service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need port number")
			}
			runServer(args[0])
			// 不会返回
			return nil
		},
	}

	return cmd
}