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
      IDP_FRONTEND_URL: process.env.IDP_FRONTEND_URL,
      IDP_SERVER_URL: process.env.IDP_SERVER_URL,
      FIFA_FRONTEND_URL: process.env.FIFA_FRONTEND_URL,
      ETH_CHAIN_ID: process.env.ETH_CHAIN_ID,
      POL_CHAIN_ID: process.env.POL_CHAIN_ID,
      BSC_CHAIN_ID: process.env.BSC_CHAIN_ID,
      RPC_URL: process.env.RPC_URL,
      ETH_SC_ADDRESS: process.env.ETH_SC_ADDRESS,
      POL_SC_ADDRESS: process.env.POL_SC_ADDRESS,
      BSC_SC_ADDRESS: process.env.BSC_SC_ADDRESS,

    }
  },
  app: {
    baseURL: process.env.NUXT_PATH_PREFIX || '',
    head: {
      link: [{ rel: "icon", type: "image/png", href: process.env.NUXT_PATH_PREFIX + "/favicon.png" }],
      title: "FIFA World Cup",
    }
  },
});
