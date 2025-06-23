import { purgeCSSPlugin } from '@fullhuman/postcss-purgecss';
import tailwindcss from "@tailwindcss/postcss"
import cssnanoPlugin from "cssnano";

export default {
  plugins: [
    tailwindcss(),

    ...(process.env.NODE_ENV === "production" ? [
      purgeCSSPlugin({
        content: [
          './**/*.html',
          './assets/**/*.css'
        ]
      }),
      cssnanoPlugin()
    ] : []),
  ]
};
