<template>
  <div v-if="!isAuthenticated">
    <SectionsNotAuthorized />
  </div>
  <div v-else class="h-full bg-stadium bg-cover">
    <div class="flex flex-col items-center py-10 bg-slate-100 bg-opacity-70 backdrop-blur-sm mx-60 rounded-lg">
    <h2 class="text-2xl font-bold">Book a Seat</h2>
    <h2 class="mt-3">Match: Argentina - France</h2>
    <div class="mt-6 flex">
      <SectionsSeat
        @zoneSelected="(z) => (selectedZone = z)"
        @seatSelected="(s) => (selectedSeat = s)"
      />
    </div>
    <BadgesSuccess
      class="mt-7"
      v-if="successMessage"
      @closeBadgeButtonClicked="() => (successMessage = false)"
    >
      Booking completed successfully
    </BadgesSuccess>
    <BadgesError
      class="mt-7"
      v-if="errorMessage"
      @closeBadgeButtonClicked="() => (errorMessage = false)"
      >{{ errorStr }}</BadgesError
    >
    <Button
      class="ml-5 mt-5 rounded bg-red-800 py-3 px-6 text-white hover:bg-red-400"
      @click="onBookButtonClick"
      v-if="!successMessage"
    >
      Book
    </Button>
    <NuxtLink
      to="/profile"
      class="ml-5 mt-5 rounded bg-red-800 py-3 px-6 text-white hover:bg-red-400"
      v-else
    >
      Go back to my profile
    </NuxtLink>
  </div>
  </div>
</template>

<script setup>
const config = useRuntimeConfig();
var errorStr;
var token;
var selectedZone;
var selectedSeat;
const isAuthenticated = ref(false);
const successMessage = ref(false);
const errorMessage = ref(false);
if (process.client) {
  token = sessionStorage.getItem("jwt");
  isAuthenticated.value = token != undefined;
}

function onBookButtonClick() {
  errorMessage.value = false;
  successMessage.value = false;
  console.log("Selected zone: " + selectedZone);
  console.log("Selected seat: " + selectedSeat);
  let bookParams = JSON.stringify({
    seat: parseInt(selectedSeat),
    zone: selectedZone,
  });
  const idpBackUrl = config.public.IDP_SERVER_URL;
  fetch(idpBackUrl + "/fifa/booking", {
    method: "POST",
    mode: "cors",
    body: bookParams,
    headers: {
      Authorization: token,
      "Content-Type": "application/json",
    },
  }).then((response) => {
    if (response.ok) {
      console.log("Book completed succesfully");
      sessionStorage.setItem("bookingInfo", bookParams);
      successMessage.value = true;
    } else {
      response.text().then((text) => {
        errorMessage.value = true;
        console.log("ERROR: " + text);
        errorStr = text;
      });
    }
  });
}
function getCookie(name) {
  if (typeof window !== "undefined") {
    var cookieArr = document.cookie.split(";");
    for (var i = 0; i < cookieArr.length; i++) {
      var cookiePair = cookieArr[i].split("=");
      if (name == cookiePair[0].trim()) {
        return decodeURIComponent(cookiePair[1]);
      }
    }
    return null;
  }
}
</script>
