"use client"
import type { Metadata } from "next";
import "../../globals.css";
import Footer from "@/app/components/footer";
import TopBar from "@/app/components/topbar";
import NewsLetter from "@/app/components/newsletter";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAccount } from "wagmi";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { isConnected } = useAccount();
  const [showWallet, setShowWallet] = useState(false);
  const [showNotifications, setShowNotifications] = useState(false);
  const [showProfile, setShowProfile] = useState(false);
  const router = useRouter();
  useEffect(() => {
    if (!isConnected) {
      router.push("/i");
    }
  }, [isConnected, router])
  if (!isConnected) {
    return <div></div>
  }
  return (
    <section onClick={() => {
      setShowWallet(false);
      setShowNotifications(false);
    }}>
      <TopBar home="/u/h" showWallet={showWallet} showNotifications={showNotifications} setShowWallet={setShowWallet} setShowNotifications={setShowNotifications} showProfile={showProfile} setShowProfile={setShowProfile} />
      {children}
      <NewsLetter />
      <Footer />
    </section>
  );
}
