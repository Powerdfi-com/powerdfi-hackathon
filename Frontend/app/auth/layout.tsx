/* eslint-disable react-hooks/exhaustive-deps */
"use client"
import Footer from "@/app/components/footer";
import { usePathname } from "next/navigation";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const path = usePathname();



  return (
    <section className="h-screen flex flex-col">
      <div className="flex-1">
        {children}
      </div>
      {
        path !== "/auth/admin" && <Footer showFull={false} />
      }
    </section>
  );
}
