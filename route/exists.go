package route

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/theTardigrade/fbdServer-v2/options"
// )

// type existsPostRequest struct {
// 	Url string `json:"u"`
// }

// type existsPostResponse struct {
// 	Success    bool `json:"s"`
// 	StatusCode int  `json:"c"`
// }

// var (
// 	existsPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var pRes existsPostResponse

// 		defer func() {
// 			// if err := recover(); err != nil {
// 			// 	serverErrorHandler(w, r)
// 			// }

// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(pRes)
// 		}()

// 		var pReq existsPostRequest

// 		err := json.NewDecoder(r.Body).Decode(&pReq)
// 		if err != nil {
// 			return
// 		}

// 		fmt.Println(pReq)

// 		res, err := http.Get(pReq.Url)
// 		if err != nil {
// 			return
// 		}

// 		pRes.Success = true
// 		pRes.StatusCode = res.StatusCode
// 	})
// )

// func init() {
// 	options.Options.Routes.Post["/exists"] = existsPostHandler
// }
