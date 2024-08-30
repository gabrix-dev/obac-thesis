<template>
  <div class="h-full bg-desert2 bg-cover py-10 text-center">
    <img src="~/assets/img/worldcupp.png" class="mx-auto w-1/6" />
    <p class="my-5 text-xl">
      Get ready for the biggest soccer event of the year!
    </p>
    <div class="flex justify-center">
      <NuxtLink
        to="/buy"
        class="rounded bg-red-800 py-3 px-10 text-white hover:bg-red-400"
        >Buy a ticket</NuxtLink
      >
      <NuxtLink
        to="/profile"
        class="ml-5 rounded bg-red-800 py-3 px-6 text-white hover:bg-red-400"
        v-if="isUserAuthenticated()"
        >Access my profile</NuxtLink
      >
      <button
        @click="onAuthenticationClick"
        class="ml-5 rounded bg-red-800 py-3 px-6 text-white hover:bg-red-400"
        v-else
      >
        Access my profile
      </button>
    </div>
    <div v-if="showErrorMsg" class="mx-20 mt-7">
      <BadgesError
        class="animate-fade mx-60"
        @closeBadgeButtonClicked="() => (showErrorMsg = false)"
      >
        {{ errorMsg }}.</BadgesError
      >
    </div>
  </div>
</template>

<script setup>
const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL + "/api/oidc/authorize";
const fifaFrontUrl = config.public.FIFA_FRONTEND_URL;

const showErrorMsg = ref(false);
var errorMsg;
var url;

onOpenHandler();

function isUserAuthenticated() {
  return process.client && sessionStorage.getItem("jwtPayload") != undefined;
}

async function onAuthenticationClick() {
  url = await createAuthParameters();
  console.log("URL: " + url);
  window.location.href = url;
}

async function createAuthParameters() {
  if (process.client) {
    const params = {
      response_type: "code",
      client_id: "abcdefgh",
      scope: "openid asset ticket wallet",
      redirect_uri: fifaFrontUrl,
    };

    // generate dynamic values and add them to the params object
    const code_verifier = generateRandomBytes(32);
    console.log(code_verifier);
    const state = generateRandomBytes(32);
    const code_challenge = await sha256(code_verifier);
    const encodedState = base64UrlEncode(state);
    params.code_challenge = code_challenge;
    params.code_challenge_method = "S256";
    params.state = encodedState;

    //store the variables
    sessionStorage.setItem("state", encodedState);
    sessionStorage.setItem("verifier", code_verifier);
    const storedState = sessionStorage.getItem("state");
    console.log(`STORED STATE: ${storedState}`);

    // create url with all parameters
    console.log("Invalid URL: " + idpBackUrl);
    const url = new URL(idpBackUrl);
    Object.keys(params).forEach((key) =>
      url.searchParams.append(key, params[key]),
    );
    return url.toString();
  } else return undefined;
}

function generateRandomBytes(n) {
  if (process.client) {
    const array = new Uint8Array(n);
    window.crypto.getRandomValues(array);
    return array;
  } else return undefined;
}

async function sha256(data) {
  console.log(data);
  const encoder = new TextEncoder();
  let dataArray = encoder.encode(data);
  console.log(dataArray);
  const hashBuffer = await crypto.subtle.digest("SHA-256", dataArray);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}

function onOpenHandler() {
  if (process.client) {
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get("error");
    if (error) {
      console.log("Error: " + error);
      errorMsg = error;
      showErrorMsg.value = true;
    } else {
      console.log("Null state or code");
    }
  }
}

function base64UrlEncode(str) {
  return btoa(
    encodeURIComponent(str).replace(
      /%([0-9A-F]{2})/g,
      function toSolidBytes(match, p1) {
        return String.fromCharCode("0x" + p1);
      },
    ),
  )
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}
</script>
