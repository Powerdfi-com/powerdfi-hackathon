"use client"
import Image from "next/image"
import { FaDiscord, FaLinkedin, FaXTwitter } from "react-icons/fa6";
import { Link } from "./link";

const Footer = ({ showFull = true }: { showFull?: boolean }) => {
    const [getStarted, features]: { text: string, href: string }[][] = [[{
        text: "Tutorials",
        href: "/"
    }, { text: "Resources", href: "/", }, { text: "Guides", href: "/" }], [{ text: "Pricing", href: "" }, { text: "Education", href: "" }, { text: "Refer a friend", href: "" }]];
    return (
        <footer className={(showFull ? "py-16" : "") + " bg-white dark:bg-primary/[0.01] backdrop-blur-lg px-6 sm:px-12 lg:px-24 ring-1 ring-primary/40 flex flex-col"} id="footer">
            {
                showFull && <section className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-12">
                    <div className="flex-1 flex" >
                        <div className="h-24 w-24 relative mx-12">
                            <Image src="/logo.png" alt="PowerDfi" fill={true} />
                        </div>
                    </div>
                    <div className="flex-1 flex flex-col gap-4">
                        <h3 className="font-semibold text-black dark:text-white">Get Started</h3>
                        <ul className="flex flex-col gap-2">
                            {getStarted.map(({ text, href: to }) => <li key={text}>
                                <Link className="text-black dark:text-white text-sm" href={to}>{text}</Link>
                            </li>)}
                        </ul>
                    </div>
                    <div className="flex-1 flex flex-col gap-4">
                        <h3 className="font-semibold text-black dark:text-white">Features</h3>
                        <ul className="flex flex-col gap-2">
                            {features.map(({ text, href: to }) => <li key={text}>
                                <Link className="text-black dark:text-white text-sm" href={to}>{text}</Link>
                            </li>)}
                        </ul>
                    </div>
                    <div className="flex-1 flex flex-col gap-4">
                        <h3 className="font-semibold text-black dark:text-white">Join the community</h3>
                        <ul className="flex gap-2">
                            <li><Link className="text-black dark:text-white" href={"/"}><FaXTwitter /></Link></li>
                            <li><Link className="text-black dark:text-white" href={"/"}><FaLinkedin /></Link></li>
                            <li><Link className="text-black dark:text-white" href={"/"}><FaDiscord /></Link></li>
                        </ul>
                        <Link className="text-black dark:text-white" href={"mailto:support@powerdfi.com"}>support@powerdfi.com</Link>
                    </div>
                </section>
            }
            <section className={(showFull ? "mt-16" : "h-[168px] items-center") + " flex gap-8 text-black dark:text-white text-sm justify-center"}>
                <div>&copy; {new Date().getFullYear()}. All rights reserved. PowerDfi.</div>
                <ul className="flex gap-8">
                    <li><Link href="/i/terms">Terms</Link></li>
                    <li><Link href="/i/privacy-policy">Privacy policy</Link></li>
                </ul>
            </section>
        </footer>
    )
}

export default Footer