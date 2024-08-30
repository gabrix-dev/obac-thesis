<template>
  <div class="h-full bg-stadium bg-cover">
    <Container>
      <div
        class="w-full rounded-xl bg-slate-100 bg-opacity-70 backdrop-blur-sm"
        v-if="isAuthenticated"
      >
        <h1></h1>
        <SectionsProfile />
        <div class="text-center">
          <button class="mt-10">
            <NuxtLink
              to="/"
              class="ml-auto mt-10 rounded-full bg-red-800 px-4 py-2 text-white hover:bg-red-400"
              @click="onLogoutButtonClick"
              >Logout</NuxtLink
            >
          </button>
        </div>
      </div>
      <div v-else-if="isAuthenticating">
        <AnimationsLoading />
      </div>
      <div v-else>
        <SectionsNotAuthorized />
      </div>
    </Container>
  </div>
</template>

<script setup>
import { Buffer } from "buffer";

definePageMeta({
  middleware: [
    function (to, from) {
      const error = to.query.error;
      const nonce = to.query.nonce;
      if (from.path == "/" || !from.path == "/book") {
        return;
      }
      if (error) {
        return navigateTo(to.path.replace("/profile", "/?error=" + error));
      }
    },
  ],
});

const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;
const isAuthenticating = ref(true);
const isAuthenticated = ref(false);
const ticket = ref(null);
onOpenHandler();

function onLogoutButtonClick() {
  sessionStorage.clear();
}

function onOpenHandler() {
  if (typeof window !== "undefined") {
    console.log("Client page opened");
    if (sessionStorage.getItem("jwt") != undefined)
      isAuthenticated.value = true;
    let searchParams = new URLSearchParams(window.location.search);
    let code = searchParams.get("code");
    let state = searchParams.get("state");
    console.log("CODE: " + code);
    if (code != undefined && state != undefined) {
      onRedirectHandler(code, state);
    }
  }
}

function onErrorHandler() {
  const urlParams = new URLSearchParams(window.location.search);
  const error = urlParams.get("error");
  window.location.href = window.location.href.replace("/profile", "/");
}

function onRedirectHandler(receivedCode, receivedState) {
  let searchParams = new URLSearchParams(window.location.search);

  isAuthenticating.value = true;
  let nonce = searchParams.get("nonce");
  console.log("Nonce: " + nonce);

  sessionStorage.setItem("nonce", nonce);

  let savedState = sessionStorage.getItem("state");
  let codeVerifier = sessionStorage.getItem("verifier");
  console.log(
    "Saved state - received state: " + savedState + " --- " + receivedState,
  );
  if (receivedState != savedState) {
    console.log("Stored and received states are different!");
  } else if (codeVerifier == null) {
    console.log("Code verifier not found!");
  } else {
    let paramsJSON = JSON.stringify({
      code: receivedCode,
      code_verifier: codeVerifier,
      nonce: nonce,
    });
    fetch(idpBackUrl + "/api/oidc/token", {
      method: "POST",
      mode: "cors",
      body: paramsJSON,
    })
      .then((response) => {
        console.log("RESPONSE RECEIVED");
        if (response.ok) {
          console.log("RESPONSE OK");
          return response.json();
        } else {
          throw new Error("Something went wrong");
        }
      })
      .then((data) => {
        console.log("TOKEN FOUNDDDD");
        console.log("DATA: " + data);
        console.log("TOKEN FOUND: " + data.ownership_token);
        sessionStorage.setItem("jwt", data.ownership_token);
        const parts = data.ownership_token.split(".");
        const decodedPayload = JSON.parse(Buffer.from(parts[1], "base64"));
        ticket.value = decodedPayload;
        sessionStorage.setItem("jwtPayload", Buffer.from(parts[1], "base64"));
        isAuthenticating.value = false;
        isAuthenticated.value = true;
        console.log(JSON.stringify(decodedPayload, null, 4));
      })
      .catch((error) => {
        console.error(error);
      });
  }
}

function getCookie(name) {
  // Split cookie string and get all individual name=value pairs in an array
  var cookieArr = document.cookie.split(";");

  // Loop through the array elements
  for (var i = 0; i < cookieArr.length; i++) {
    var cookiePair = cookieArr[i].split("=");

    /* Removing whitespace at the beginning of the cookie name and compare it with the given string */
    if (name == cookiePair[0].trim()) {
      // Decode the cookie value and return
      return decodeURIComponent(cookiePair[1]);
    }
  }

  // Return null if not found
  return null;
}
</script>

<script>
export default {
  data() {
    return {
      ticket: {
        match: "Argentina - France",
        seatType: "Normal",
        zones: "CAT1, CAT2",
        dinner: "Snack + drink",
        image: "/_nuxt/assets/img/fifa_basic_ticket.png",
      },
    };
  },
};
</script>
