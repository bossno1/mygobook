package main

import(
    "fmt"
    "crypto/sha256"
   
)

func main() {

    // 第一种调用方法
    sum := sha256.Sum256([]byte("app_id=3vtzutuzb131il4io4&biz_content=c30d18ac98e71c03eaa110aab94b5abed6439c942325d1ab50d60e7d8a8460844c9eaf96b66083445bb3f030c3b3782f9f88e382dafc5091f7b2ad63e6994b9a368d52ee36c8dbd7429427d8ca6b6a8a27d4f23faddf49f5600fa340c9a2376d&digest_type=SM3&enc_type=SM4&method=ehc.ehealthcode.verify&term_id=SK0001&timestamp=1557236049417&version=X.M.0.1ae4ead79e2bb40d69f163717b61e7984"))
    fmt.Printf("%x\n", sum)
	
    // 第二种调用方法
    h := sha256.New()
    h.Write([]byte("app_id=3vtzutuzb131il4io4&biz_content=wxDt0eO7oP6xpT+486dltG29VynpbkZowWu5STu/9avVPywIIGbz3coJTcQ/31GAUjbM0iXlPPldHz8339dW5iCibd37TIy5mxyV9TyNc6knQGeY/fyNAM77+NjFvY+pCFHEJIy7z4zbYdQxhkMNOwjteIqtsa78h/NkIzstTM72UR2wrZANYNPonZUJhPTMwIGuZUu5MAu59sYRZncsj31022t8YkTDPWXhbxhNArCEnnzk7KiEkIW4au8xVs2eATwqdOChZ6625yQmQiCFyk5TVBGvWqpkZi/GZGtdE3krPAnqhQG9SnOk+fZLVLgjwNyBWw7DSE4qblXrf1xkng==ae4ead79e2bb40d69f163717b61e7984&digest_type=SHA256&enc_type=AES&method=ehc.ehealthcard.register&term_id=&timestamp=1558339096000&version=X.M.0.1ae4ead79e2bb40d69f163717b61e7984"))
    fmt.Printf("%x\n", h.Sum(nil))
}