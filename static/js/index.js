document.addEventListener("DOMContentLoaded", function () {
  var seeMoreButtons = document.querySelectorAll(".seeMoreButton");

  seeMoreButtons.forEach(function (button) {
    button.addEventListener("click", function () {
      var commentsContainer = button.parentElement.querySelector(".comments");
      
      if (commentsContainer.style.maxHeight === "110px" || commentsContainer.style.maxHeight === "") {
        commentsContainer.style.maxHeight = "100%"; 
        button.textContent = "See Less";
      } else {
        commentsContainer.style.maxHeight = "110px";
        button.textContent = "See More";
      }
    });
  });
});