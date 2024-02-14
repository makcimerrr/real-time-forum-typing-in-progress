import { getCookie } from "./cookie.js";
import { showDiv } from "./show.js";
import { login } from "./login.js";
import { post, fetchAndDisplayDiscussions } from "./post.js";
import {register} from "./register.js";
import {startWebSocket} from "./websocket.js";
import {logout} from "./websocket.js";

document.addEventListener('DOMContentLoaded', function () {
    // Appeler la fonction pour récupérer et afficher les discussions initiales
    fetchAndDisplayDiscussions();
    // Gérer le clic sur le titre de chaque discussion pour afficher le message
});

window.onload = function () {
    const username = getCookie("username");
    if (username) {
        startWebSocket()
        showDiv("home");
    }
}

document
    .getElementById("Forum")
    .addEventListener("click", function () {
        showDiv("home");
    });

document
    .getElementById("logout")
    .addEventListener("click", function (event) {
        // Empêcher le comportement de soumission par défaut
        event.preventDefault();

        // Appeler la fonction login de votre script JavaScript
        logout();
    });

document
    .getElementById("loginForm")
    .addEventListener("submit", function (event) {
        // Empêcher le comportement de soumission par défaut
        event.preventDefault();

        // Appeler la fonction login de votre script JavaScript
        login();
    });

document
    .getElementById("registrationForm")
    .addEventListener("submit", function (event) {
        // Empêcher le comportement de soumission par défaut
        event.preventDefault();

        // Appeler la fonction login de votre script JavaScript
        register();
    });

document
    .getElementById("postForm")
    .addEventListener("submit", function (event) {
        // Empêcher le comportement de soumission par défaut
        event.preventDefault();

        // Appeler la fonction login de votre script JavaScript
        post();
    });