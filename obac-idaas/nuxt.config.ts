export default defineNuxtConfig({
  css: [
    "~/assets/css/main.css",
    "@fortawesome/fontawesome-svg-core/styles.css",
  ],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  runtimeConfig: {
    public: {
      IDP_SERVER_URL: process.env.IDP_SERVER_URL,
      RPC_URL: process.env.RPC_URL,
      ETH_CHAIN_ID: process.env.ETH_CHAIN_ID,
      POL_CHAIN_ID: process.env.POL_CHAIN_ID,
      BSC_CHAIN_ID: process.env.BSC_CHAIN_ID,
    }
  },
  app: {
    baseURL: process.env.NUXT_PATH_PREFIX || '',
    head: {
      link: [{ rel: "icon", type: "image/png", href: process.env.NUXT_PATH_PREFIX + "/favicon.png" }],
      title: "OBAC IDaaS",
    }
  }
});
