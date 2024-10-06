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
    return (
      <>
        <title>{title} – GIMBAP</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta property="og:title" content={title + " – GIMBAP"} />
        <meta
          property="og:description"
          content="Go Injection Management for Better Application Programmingr"
        />
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
