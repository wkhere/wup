package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 9000, "")
	flag.Parse()

	fmt.Fprintf(os.Stderr, uploadInfo(*port))

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func uploadInfo(port int) string {
	return fmt.Sprintf(`
to upload, run:  curl --data-binary @src http://localhost:%d/dest
            or:  httpstat -d @src http://localhost:%d/dest
or ssh -R $remoteport:localhost:%d $host,
    then on $host: curl .... http://localhost:$remoteport/dest
or map your revproxy /wup to localhost:%d and use http://$host/wup/dest
then find your data in %s/dest

`,
		port, port, port, port, destDir,
	)
}
