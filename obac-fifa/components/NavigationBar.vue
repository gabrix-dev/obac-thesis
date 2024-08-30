<template>
  <header class="flex items-center justify-between bg-red-800 p-4 text-white">
    <!-- <IconsBall /> -->
    <!-- <IconsWorldCup class="h-11 w-11" /> -->
    <div class="flex">
      <img src="~/assets/img/ball.svg" class="mr-5 w-8 fill-white" />
      <!-- <IconsBall class="mr-5 w-8" /> -->
      <h1 class="text-xl font-bold">Fifa world cup</h1>
    </div>

    <nav>
      <ul class="flex">
        <li class="px-5">
          <NuxtLink to="/" class="text-white hover:text-gray-400"
            >Home</NuxtLink
          >
        </li>
        <li class="px-5">
          <NuxtLink to="/buy" class="text-white hover:text-gray-400"
            >Buy a ticket</NuxtLink
          >
        </li>
      </ul>
    </nav>
  </header>
</template>
