package httpserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func HTTPServer(baseOutPath, chunkListFilename, serveHttpAddr string, log *logrus.Logger) {
	fs := http.FileServer(http.Dir(baseOutPath))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html lang="en">
				<head>
					<meta charset=utf-8/>
					<script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
				</head>
				<body>
				<video
					id="my-player"
					autoplay
					loop
					muted="muted"
					controls
					style="border: 5px solid #000"
					data-setup='{}'>
				</video>
				<script>
					var link = "/video/`+ chunkListFilename + `";
					var video = document.querySelector('#my-player');
					video.src = link;

					if (Hls.isSupported()) {
					var hls = new Hls();
					hls.loadSource(link);
					hls.attachMedia(video);
					hls.on(Hls.Events.MANIFEST_PARSED, function () {
						video.play();
					});
					}
					else if (video.canPlayType('application/vnd.apple.mpegurl')) {
					video.src = link;
					video.addEventListener('loadedmetadata', function () {
						video.play();
					});
					}
				</script>
				</body>
			</html>`))
	})
	http.Handle("/video/", NoCache(http.StripPrefix("/video/", fs)))

	go http.ListenAndServe(serveHttpAddr, nil)

	log.Printf("HTTP server listening on %s", serveHttpAddr)
}
