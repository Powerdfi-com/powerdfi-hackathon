import { Metadata } from "next";
import Image from "next/image";
import { IoPersonSharp } from "react-icons/io5";
import GhostBlogs from "./blogs";

export const metadata: Metadata = {
  title: "PowerDfi - Blogs",
  description: "Stay updated with the latest insights, trends, and news on how blockchain technology is revolutionizing the management and trading of real- world assets (RWAs).",
};



export default function Blog() {


  return (
    <main className="bg-gradient-radial from-black to-primary/5">

      <GhostBlogs />
    </main>
  );
}
