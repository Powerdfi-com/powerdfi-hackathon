"use client"
import { MdLightMode } from "react-icons/md";
import { MdDarkMode } from "react-icons/md";


import Image from "next/image"
import { Link } from "./link";
import { useAccount } from "wagmi";
import { isDarkMode, switchToDarkMode, switchToLightMode } from "../utils/func";
import { useContext, useEffect, useState } from "react";
import { MediaQueryContext, ThemeContext } from "../context/context";
import { HiMenuAlt3 } from "react-icons/hi";
const TopBar = () => {
    const navLinks = [{
        text: "Service",
        href: "#services"
    }, {
        text: "About",
        href: "/i/about-us"
    }, {
        text: "Features",
        href: "#features"
    }, {
        text: "Contact",
        href: "#footer"
    }, {
        text: "Blog",
        href: "/i/blogs"
    }
    ];
    const isMobile = useContext(MediaQueryContext);
    const [showNav, setShowNav] = useState(false)

    const { isConnected } = useAccount();
    const { isDarkMode, setDarkMode } = useContext(ThemeContext);
    if (isMobile) {
        return <section className="h-16 relative flex justify-between px-8 items-center">
            <Link href={"/i"}><Image src="/logo.png" alt="logo" height={50} width={50} /></Link>
            <div className="h-10 cursor-pointer w-10 aspect-square bg-secondary flex items-center justify-center rounded-lg" onClick={() => setShowNav(!showNav)}>
                <HiMenuAlt3 className="!text-white !text-2xl" />
            </div>
            {
                (showNav && isMobile) && <div className="fixed z-20 top-0 left-0 right-0 bottom-0 bg-black/20 backdrop-blur-lg flex items-end" onClick={() => setShowNav(false)}>
                    <div className="flex w-full bg-black/80 p-6" onClick={(e) => e.stopPropagation()}>
                        <nav className="w-full">
                            <ul className="flex flex-col gap-2 w-full">
                                {navLinks.map(({ text, href }) => <li key={text} className="ring-1 ring-border w-full h-12 flex items-center rounded-md"><Link href={href} className="text-white text-md font-semibold !p-5 w-full">{text}</Link></li>)}
                            </ul>
                        </nav>
                    </div>
                </div>
            }
        </section>
    }

    return (
        <section className="flex items-center px-8 sm:px-12 md:px-24 h-[89px] gap-12 bg-text-grey">
            <Link href="/i">
                <Image src="/logo.png" alt="logo" height={50} width={50} />
            </Link>
            <nav className="flex-1 flex justify-center">
                <ul className="flex gap-12">
                    {navLinks.map(({ text, href }) => <li key={text}><Link href={href} className="text-white text-[16px] font-semibold">{text}</Link></li>)}
                </ul>
            </nav>
            <div className="flex gap-6 items-center">
                {
                    isDarkMode ? <MdLightMode onClick={() => {
                        switchToLightMode();
                        setDarkMode(false);
                    }} className="!text-4xl cursor-pointer text-secondary" /> : <MdDarkMode onClick={() => {
                        switchToDarkMode();
                        setDarkMode(true);
                    }} className="!text-4xl cursor-pointer !text-secondary" />
                }
                <Link href={isConnected ? "/u/h" : "/auth"} className="px-6 py-2 rounded-md bg-secondary text-sm">
                    Join PowerDfi Studio
                </Link>
            </div>
        </section>
    )
}

export default TopBar   