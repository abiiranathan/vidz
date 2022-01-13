import "video.js/dist/video-js.css";
import videojs from "video.js";
import "../src/index.css";
import playlist from "./videojs-playlist";

const player = videojs("#videojs", {
  loop: false,
  autoplay: true,
  controls: true,
  responsive: true,
  aspectRatio: "16:9",
  fill: true,
  playbackRates: [0.7, 1.0, 1.5, 2.0, 2.5, 3.0],
  defaultVolume: 1,
});

videojs.registerPlugin("playlist", playlist);

function createSource(video) {
  return {
    name: video.title,
    sources: [
      {
        src: "/media" + video.path,
        type: video.type,
      },
    ],
  };
}

async function fetchAllVideos(url) {
  const response = await fetch(url);
  const videos = await response.json();
  return videos;
}

async function getPlaylist() {
  const videos = await fetchAllVideos("/api/videos");

  const playlist = [];
  for (const video of videos) {
    playlist.push(createSource(video));
  }
  return playlist;
}

async function initPlaylist() {
  const playlist = await getPlaylist();
  player.playlist(player, playlist);
}

function initSearch() {
  const search_input = document.getElementById("search_input");

  search_input.addEventListener("keydown", function (e) {
    e.stopImmediatePropagation();
  });

  search_input.addEventListener("keyup", function (event) {
    if (event.key === "Enter") {
      const search_value = search_input.value;
      const search_url = "/api/videos?title=" + search_value;

      fetchAllVideos(search_url).then(videos => {
        const playlist = [];
        for (const video of videos) {
          playlist.push(createSource(video));
        }
        player.playlist(player, playlist);
      });
    }
  });
}

initPlaylist();
initSearch();
