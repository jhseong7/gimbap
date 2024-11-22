import React from "react";
import { DocsThemeConfig, useConfig } from "nextra-theme-docs";

const config: DocsThemeConfig = {
  logo: (
    <span
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        gap: "1rem",
        fontSize: "1.5rem",
        fontWeight: "bold",
      }}
    >
      <img src="/gimbap-logo.png" style={{ width: "2.5rem" }} />
      <b>GIMBAP</b>
      <p style={{ fontSize: "1rem" }}>
        Go Injection Management for Better Application Programming
      </p>
    </span>
  ),
  head: () => {
    const { title } = useConfig();

    const titleStr = title + " – GIMBAP";
    const description =
      "Go Injection Management for Better Application Programming";
    const siteRoot = "https://www.go-gimbap.com";

    return (
      <>
        <title>{titleStr}</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="description" content={description} />
        <meta property="og:title" content={titleStr} />
        <meta property="og:description" content={description} />
        <meta property="og:url" content={siteRoot} />
        <meta property="og:type" content="website" />
        <meta property="og:image" content={`${siteRoot}/og.png`} />
        <meta property="og:site_name" content={titleStr} />
        <meta property="og:logo" content={`${siteRoot}/og.png`} />
        <meta property="og:image:alt" content="GIMBAP" />
        <meta name="twitter:image" content={`${siteRoot}/og.png`} />
        <meta property="twitter:title" content={titleStr} />
        <meta property="twitter:description" content={description} />
        <meta property="twitter:card" content={`${siteRoot}/og.png`} />
        <meta name="twitter:image:type" content="website" />
        <meta name="twitter:image:width" content="800" />
        <meta name="twitter:image:height" content="420" />
        <meta property="twitter:image:alt" content="GIMBAP" />
      </>
    );
  },
  navigation: {
    prev: true,
    next: true,
  },
  editLink: {
    component: null,
  },
  feedback: {
    content: null,
  },
  footer: {
    content: `MIT 2024 ©jhseong7`,
  },
  project: {
    link: "https://github.com/jhseong7/gimbap",
    // icon: <img src="/gimbap-logo.png" style={{ width: "2.5rem" }} />,
  },
  docsRepositoryBase: "https://github.com/jhseong7/gimbap/tree/main/docs",
};

export default config;
