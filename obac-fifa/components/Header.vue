<template>
  <!-- DESKTOP -->
  <Disclosure
    id="app"
    as="nav"
    class="sticky top-0 z-30 hidden sm:block"
    v-show="showNavbar"
  >
    <div
      class="flex h-20 items-center justify-between shadow-xl backdrop-blur-sm"
    >
      <!-- Logo -->
      <a class="ml-16" href="/">
        <LogosSmallLogo />
      </a>

      <!-- PC Nav -->
      <div class="mr-16 flex">
        <div class="hidden sm:ml-6 sm:block">
          <div class="flex space-x-4">
            <TitlesHeaderTitle
              v-for="item in navigation"
              :sectionNumber="item.id"
              :key="item.name"
              :hrefs="item.href"
              :sectionName="item.name"
              :class="[
                item.current
                  ? 'bg-gray-900 text-red-800'
                  : 'text-red-800 hover:bg-gray-700 hover:text-white',
                'rounded-md px-3 py-2 font-orbitron text-base font-medium lg:text-lg',
              ]"
            >
            </TitlesHeaderTitle>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>

  <!-- MOBILE -->
  <Disclosure id="app" as="nav" class="block sm:hidden" v-slot="{ open }">
    <div class="relative flex h-16 items-center justify-between py-14 px-3">
      <!-- Logo -->
      <a href="/">
        <LogosSmallLogo />
      </a>

      <!-- Mobile menu button-->
      <div class="flex items-center sm:hidden">
        <DisclosureButton
          class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white hover:bg-gray-700 hover:text-white"
        >
          <span class="sr-only">Open main menu</span>
          <Bars3Icon
            v-if="!open"
            class="block h-8 w-8 text-neon-green"
            aria-hidden="true"
          />
        </DisclosureButton>
      </div>
    </div>

    <!-- Mobile Nav -->
    <DisclosurePanel class="sm:hidden">
      <div
        class="fixed right-0 top-0 z-10 h-screen w-4/5 rounded-xl bg-gray-800 px-2 pt-2 pb-3 opacity-95"
      >
        <DisclosureButton
          class="float-right mt-8 mr-2.5 rounded-md p-2 text-gray-400 ring-2 ring-gray-200 focus:ring-inset hover:bg-gray-700 hover:text-white"
        >
          <XMarkIcon class="h-7 w-7 text-neon-green" aria-hidden="true" />
        </DisclosureButton>
        <div class="mx-20 mt-32 text-justify">
          <DisclosureButton
            v-for="item in navigation"
            :key="item.name"
            as="a"
            :href="item.href"
            :class="[
              item.current
                ? 'bg-gray-900 text-white'
                : 'text-gray-300 hover:bg-gray-700 hover:text-white',
              'block rounded-md py-4 font-orbitron text-lg font-medium',
            ]"
            ><p class="font-orbitron text-black">{{ item.id }}.</p>
            {{ item.name }}</DisclosureButton
          >
          <p class="mt-24 text-justify text-lg text-gray-300">
            Made with <IconsHeart class="inline" /> by
          </p>
          <p class="mt-1.5 text-justify text-lg text-gray-300">
            Albert Ferraté
          </p>
        </div>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>

<script setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from "@headlessui/vue";
import {
  Bars3Icon,
  BellIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline/index.js";

const navigation = [
  { id: "01", name: "Buy Tickets", href: "#buytickets", current: false },
  { id: "02", name: "Book a sit", href: "#bookasit", current: false },
];
</script>

<script>
export default {
  data() {
    return {
      showNavbar: true,
      lastScrollPosition: 0,
    };
  },
  mounted() {
    window.addEventListener("scroll", this.onScroll);
  },
  beforeDestroy() {
    window.removeEventListener("scroll", this.onScroll);
  },
  methods: {
    onScroll() {
      const currentScrollPosition =
        window.pageYOffset || document.documentElement.scrollTop;
      if (currentScrollPosition < 0) {
        return;
      }
      // Stop executing this function if the difference between
      // current scroll position and last scroll position is less than some offset
      if (Math.abs(currentScrollPosition - this.lastScrollPosition) < 100) {
        return;
      }
      this.showNavbar = currentScrollPosition < this.lastScrollPosition;
      this.lastScrollPosition = currentScrollPosition;
    },
  },
};
</script>
