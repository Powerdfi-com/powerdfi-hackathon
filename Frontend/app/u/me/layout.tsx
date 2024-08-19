"use client"
import { GoArrowSwitch, GoHomeFill, GoVerified } from "react-icons/go";
import Footer from "@/app/components/footer";
import TopBar from "@/app/components/topbar";
import { FaRegEdit } from "react-icons/fa";
import { IoGridOutline, IoWalletOutline } from "react-icons/io5";
import { CgNotifications, CgProfile } from "react-icons/cg";
import { BiSupport } from "react-icons/bi";
import { useAccount } from "wagmi";
import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Link } from "@/app/components/link";
import { isLinkActive } from "@/app/utils/func";

export default function MeLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {

    const [showWallet, setShowWallet] = useState(false);

    const [showProfile, setShowProfile] = useState(false);
    const path = usePathname();
    const [showNotifications, setShowNotifications] = useState(false);
    const marketplace = [
        {
            title: "Overview",
            href: "/u/me",
            Icon: GoHomeFill,
        },
        {
            title: "Create",
            href: "/u/me/create",
            Icon: FaRegEdit,
        }, {
            title: "Listing",
            href: "/u/me/listing",
            Icon: GoArrowSwitch,
        }, {
            title: "KYC",
            href: "/u/me/kyc",
            Icon: GoVerified,
        }, {
            title: "Portfolio",
            href: "/u/me/portfolio",
            Icon: GoArrowSwitch,
        }, {
            title: "Explore",
            href: "/u/h/explore",
            Icon: IoGridOutline,
        }
    ]
    const settings = [
        {
            title: "Profile",
            href: "/u/me/profile",
            Icon: CgProfile,
        },
        {
            title: "Wallet",
            href: "/u/me/wallet",
            Icon: IoWalletOutline,
        }, {
            title: "Notifications",
            href: "/u/me/notifications",
            Icon: CgNotifications,
        }, {
            title: "Help Center",
            href: "/u/me/help-center",
            Icon: BiSupport,
        },
    ]
    const { isConnected } = useAccount();
    const router = useRouter();
    useEffect(() => {
        if (!isConnected) {
            router.push("/i");
        }
    }, [isConnected])
    if (!isConnected) {
        return <div></div>
    }
    return (
        <section className="h-screen flex flex-col">
            <TopBar showProfile={showProfile} setShowProfile={setShowProfile} home="/u/me" showWallet={showWallet} showNotifications={showNotifications} setShowWallet={setShowWallet} setShowNotifications={setShowNotifications} />
            <div className="flex-1 flex w-full h-full px-24">
                <section className="w-72 flex flex-col overflow-auto gap-6  max-h-[calc(100vh-280px)] pb-8 border-r border-secondary">
                    <ul className="gap-6 flex flex-col">
                        <h4 className="text-secondary text-[24px] font-semibold">MARKETPLACE</h4>
                        {
                            marketplace.map(({ title, href, Icon }, index) => <Link href={href} key={title} className={(isLinkActive(path, href, index === 0) ? "text-white" : "text-white/40") + " flex gap-4 text-[24px] items-center text-white"}>
                                <Icon />
                                <span className="text-sm">{title}</span>
                            </Link>)
                        }
                    </ul>
                    <ul className="gap-6 flex flex-col">
                        <h4 className="text-secondary text-[24px] font-semibold">SETTINGS</h4>
                        {
                            settings.map(({ title, href, Icon }) => <Link href={href} key={title} className={(isLinkActive(path, href) ? "text-white" : "text-white/40") + " flex gap-4 text-[24px] items-center text-white"}>
                                <Icon />
                                <span className="text-sm">{title}</span>
                            </Link>)
                        }
                    </ul>
                </section>
                <section className="flex-1 overflow-auto px-4 max-h-[calc(100vh-320px)]">
                    {children}
                </section>
            </div>
            <Footer showFull={false} />
        </section>
    );
}
