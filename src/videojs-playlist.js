let initialised = false;

function playlist(player, vid_playlist) {
  const playlistContainer = document.getElementById("vjs-playlist");
  const videoContainer = player.el();

  playlistContainer.innerHTML = "";
  let currentIndex = 0;

  vid_playlist.forEach(createPlaylistItem);
  setNowPlaying(vid_playlist[0]);

  function handleKeyDown(e) {
    if (e.code === "Space") {
      e.preventDefault();
      e.stopPropagation();

      if (player.paused()) {
        player.play();
      } else {
        player.pause();
      }
    } else if (e.code === "KeyN") {
      playNextSource();
    } else if (e.code === "KeyP") {
      playPrevSource();
    }
  }

  if (!initialised) {
    window.addEventListener("keydown", handleKeyDown);
    player.on("ended", playNextSource);
    player.on("loadedmetadata", highlightActive);

    initialised = true;
  }

  function playNextSource() {
    const video = vid_playlist[currentIndex + 1];
    if (video) {
      currentIndex += 1;
      setNowPlaying(video);
    }
  }

  function playPrevSource() {
    const video = vid_playlist[currentIndex - 1];
    if (video) {
      currentIndex -= 1;
      setNowPlaying(video);
    }
  }

  function highlightActive() {
    const allTracks = document.querySelectorAll(".vjs-playlist-item");
    allTracks.forEach(tr => {
      tr.classList.remove("playing");
    });

    const activeTrack = document.getElementById(currentIndex.toString());

    if (activeTrack) {
      activeTrack.classList.add("playing");

      if (window.scrollY < 43) {
        // Video is in full view
        return;
      }

      // scroll video into view
      videoContainer.scrollIntoView(true, {
        behavior: "smooth",
      });
    }
  }

  function createPlaylistItem(video, index) {
    const li = document.createElement("li");
    li.className = "vjs-playlist-item";
    li.id = index.toString();
    li.textContent = video.name;

    li.addEventListener("click", () => {
      currentIndex = index;
      setNowPlaying(video);
    });

    playlistContainer.appendChild(li);
  }

  function setNowPlaying(video) {
    if (!video) return;
    player.src(video.sources);
  }
}

export default playlist;
