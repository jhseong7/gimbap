import type { AppProps } from "next/app";

import { Analytics } from "@vercel/analytics/next";
import { SpeedInsights } from "@vercel/speed-insights/next";

export default function App({ Component, pageProps, router }: AppProps) {
  return (
    <>
      <Component key={router.route} {...pageProps} />
      <Analytics />
      <SpeedInsights />
    </>
  );
}
