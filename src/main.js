import "video.js/dist/video-js.css";
import VideJS from "video.js";
const element = document.querySelector("#videojs");

new VideJS(element, {
  loop: true,
  autoplay: true,
  controls: true,
  responsive: true,
  aspectRatio: "16:9",
  fill: true,
});

const deleteButton = document.querySelector(".btn-delete");
deleteButton.addEventListener("click", e => {
  e.preventDefault();

  const id = deleteButton.getAttribute("data-id");
  handledelete(id);
});

function handledelete(id) {
  const url = "/?id=" + id;

  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then(res => {
      if (res.status === 200) {
        return res.json();
      }

      throw new Error("Request failed!");
    })
    .then(() => {
      window.location.href = "/";
    })
    .catch(err => {
      console.log(err);
    });
}
