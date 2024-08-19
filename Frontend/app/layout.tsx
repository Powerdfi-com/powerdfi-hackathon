import type { Metadata } from "next";
import "./globals.css";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import 'react-toastify/dist/ReactToastify.css';
import 'nprogress/nprogress.css';
import { Sofia_Sans } from "next/font/google";
import Provider from "./provider";
const sophia = Sofia_Sans({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "PowerDfi",
  description: "We aim to revolutionise the way real-world assets (RWAs) are managed and traded by leveraging blockchain technology.",
  icons: [
    {
      rel: "icon",
      type: "image/x-icon",
      url: "/icon.png",
    },
    {
      rel: "apple-touch-icon",
      sizes: "180x180",
      url: "/icon.png",
    }
  ]
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html className={sophia.className + " bg-background scroll-smooth"} lang="en">
      <head>
        <title>PowerDfi</title>
      </head>
      <body>
        <Provider>
          {children}
        </Provider>
      </body>
    </html>
  );
}
