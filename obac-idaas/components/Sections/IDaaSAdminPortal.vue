<template>
  <div class="flex h-16 items-center justify-end bg-entrust-purple px-16">
    <span class="text-xl font-semibold text-white">FIFA Admin User</span>
    <a href="#" class="mx-4 rounded-full bg-white py-1 px-2">Logout</a>
  </div>
  <div
    class="absolute flex h-screen flex-grow flex-col overflow-y-auto border-r border-gray-200 bg-white pt-5 pb-4"
  >
    <div class="flex flex-shrink-0 items-center px-4">
      <LogosIDaaS class="h-12 w-56" />
    </div>
    <div class="mt-5 flex flex-grow flex-col">
      <nav class="flex-1 space-y-1 bg-white pl-2 pr-3" aria-label="Sidebar">
        <a
          v-for="item in navigation"
          :key="item.name"
          :href="item.href"
          :class="[
            item.current
              ? 'bg-gray-100 text-gray-900'
              : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
            'group flex items-center rounded-md px-2 py-2 text-sm font-medium',
          ]"
        >
          <component
            :is="item.icon"
            :class="[
              item.current
                ? 'text-gray-500'
                : 'text-gray-400 group-hover:text-gray-500',
              'mr-3 h-6 w-6 flex-shrink-0',
            ]"
            aria-hidden="true"
          />
          <span class="flex-1">{{ item.name }}</span>
          <span
            v-if="item.count"
            :class="[
              item.current ? 'bg-white' : 'bg-gray-100 group-hover:bg-gray-200',
              'ml-3 inline-block rounded-full py-0.5 px-3 text-xs font-medium',
            ]"
            >{{ item.count }}</span
          >
        </a>
      </nav>
    </div>
  </div>
  <div class="flex">
    <div class="w-64">
      <!-- We fake the absolute div in order to center the rest of the content -->
    </div>
    <div class="mt-12 flex w-full justify-center">
      <div
        class="w-2/3 justify-center rounded-xl border-t-2 bg-white shadow-lg sm:w-4/5 2xl:w-2/3"
      >
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
// import { MapIcon } from "@heroicons/vue/20/solid";
import {
  ChartBarIcon,
  ClipboardDocumentCheckIcon,
  FolderIcon,
  HomeIcon,
  InboxIcon,
  MapIcon,
} from "@heroicons/vue/24/outline";

const props = defineProps({
  currentPage: {
    type: String,
    required: false,
  },
});

const navigation = [
  { name: "Dashboard", icon: HomeIcon, href: "#", current: false },
  { name: "Identity Policies", icon: FolderIcon, href: "#", current: false },
  {
    name: "Ownership Policies",
    icon: FolderIcon,
    href: "/zt/dev/obac-idaas/policy",
    current: props.currentPage == "policy",
  },
  {
    name: "Ownership Claims ",
    icon: MapIcon,
    href: "/zt/dev/obac-idaas/oidc",
    current: props.currentPage == "oidc",
  },
  {
    name: "Whitelist DBs",
    icon: ClipboardDocumentCheckIcon,
    href: "#",
    current: false,
  },
  { name: "Customization", icon: ChartBarIcon, href: "#", current: false },
  { name: "Reports", icon: InboxIcon, href: "#", count: 12, current: false },
];
</script>
