<template>
  <div>
    <div class="mx-8 flex items-center justify-between">
      <img src="~/assets/img/Argentina_flag_icon.svg" class="w-32" />
      <h2 class="mt-7 w-full text-center text-2xl font-bold text-black">
        Hi ticket owner, welcome to your profile
      </h2>
      <img src="~/assets/img/france_icon.png" class="ml-5 w-32" />
    </div>
    <SectionsTicketInfo />
    <div v-if="isLoadingBookingInfo">
      <AnimationsLoading />
    </div>
    <div v-else>
      <div v-if="!isTicketBooked" class="text-center">
        <img :src="imageUri" class="mx-auto w-4/12 py-4" />
        <h2 class="pb-6">You haven't booked a seat yet!</h2>
        <NuxtLink
          to="/book"
          class="ml-5 rounded bg-red-800 py-3 px-10 text-white hover:bg-red-400"
          >Book a seat now</NuxtLink
        >
      </div>
      <div v-else class="mt-10 mb-5 flex">
        <Booking :bookingInfoObject="bookingInfoObject" />
      </div>
    </div>
  </div>
</template>

<script setup>
import fifaBasicTicketImg from "~/assets/img/fifa_basic_ticket-transparent.png";
import fifaBasicPlusTicketImg from "~/assets/img/fifa_basicplus_ticket-transparent.png";
import fifaVipTicketImg from "~/assets/img/fifa_vip_ticket-transparent.png";

const config = useRuntimeConfig();
const isTicketBooked = ref(false);
const isLoadingBookingInfo = ref(true);
var bookingInfoObject;
var token;
if (process.client) {
  token = sessionStorage.getItem("jwt");
}
getBookingInfo();

async function getBookingInfo() {
  const idpBackUrl = config.public.IDP_SERVER_URL;
  fetch(idpBackUrl + "/fifa/booking", {
    method: "GET",
    mode: "cors",
    headers: {
      Authorization: token,
    },
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("The user hasn't booked a seat yet");
      }
    })
    .then((bookingInfo) => {
      let str = JSON.stringify(bookingInfo, null, 4);
      console.log("Booking registry stringified: " + str);
      bookingInfoObject = bookingInfo;
      isLoadingBookingInfo.value = false;
      isTicketBooked.value = true;
    })
    .catch((error) => {
      isLoadingBookingInfo.value = false;
      isTicketBooked.value = false;
    });
}

var ticket;

if (process.client) {
  let ticketStr = sessionStorage.getItem("jwtPayload");
  ticket = JSON.parse(ticketStr);
}

var imageUri = getImage();

function getImage() {
  if (process.client) {
    console.log(ticket.ticket_type);
    switch (ticket.ticket_type) {
      case "Basic":
        return fifaBasicTicketImg;
      case "Basic plus":
        return fifaBasicPlusTicketImg;
      default:
        return fifaVipTicketImg;
    }
  } else {
    return "";
  }
}
</script>
