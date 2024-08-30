<template>
  <div class="py-5 px-10">
    <h2 class="py-4 text-center text-lg font-bold underline">
      Your ticket info
    </h2>
    <div class="flex justify-center text-black">
      <ul class="w-max space-y-4 text-left">
        <!-- <li class="text-lg font-bold">Match: {{ ticket.match }}</li> -->
        <li class="text-lg font-bold">Seat type: {{ ticket.seat_type }}</li>
        <li class="text-lg font-bold">Bookable zones: {{ ticket.zones }}</li>
        <li class="text-lg font-bold">Dinner: {{ ticket.dinner }}</li>
      </ul>
    </div>
    <!-- <img :src="imageUri" class="mx-auto w-4/12 py-4" /> -->
  </div>
</template>

<script setup>
var ticket;

if (process.client) {
  let ticketStr = sessionStorage.getItem("jwtPayload");
  ticket = JSON.parse(ticketStr);
}
///_nuxt/assets/img/worldcupp.pn
// var baseImgUri = "~/assets/img/";
var baseImgUri = "/_nuxt/assets/img/";
var imageUri = baseImgUri + getImageName();

function getImageName() {
  if (process.client) {
    console.log(ticket.ticket_type);
    switch (ticket.ticket_type) {
      case "Basic":
        return "fifa_basic_ticket.png";
      case "Basic plus":
        return "fifa_basicplus_ticket.png";
      default:
        return "fifa_vip_ticket.png";
    }
  } else {
    return "";
  }
}

// const props = defineProps({
//   ticket: {
//     type: Object,
//     required: true,
//   },
// });
</script>
