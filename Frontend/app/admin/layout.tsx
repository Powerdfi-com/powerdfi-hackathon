"use client"
import "../globals.css"
import Footer from "@/app/components/footer";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import TopBar from "@/app/components/landingtopbar";
import { Link } from "../components/link";
import { isLinkActive } from "../utils/func";
import { usePathname, useRouter } from "next/navigation"; import { IoSettingsOutline } from "react-icons/io5"; import { CiLogout } from "react-icons/ci";
import { MdDashboard } from "react-icons/md"; import { IoPersonCircleOutline } from "react-icons/io5"; import { CiCircleInfo } from "react-icons/ci";
import { useCallback, useEffect, useState } from "react";
import { toast } from "react-toastify";
import { cookie } from "../utils/cookie";
import AuthAPI from "../utils/apis/auth";
import { TCredentials } from "../utils/types";
import { IoIosArrowForward } from "react-icons/io";
import Modal from "../components/modal";



export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const path = usePathname();
    const router = useRouter();
    const [showLogout, setShowLogout] = useState(false);
    const [showMore, setShowMore] = useState(false);
    const [credentials, setCredentials] = useState<TCredentials>(cookie.getJson("adminCredentials"));
    const signOut = useCallback(() => {
        toast.error("Credentials expired, please sign in again!")
        cookie.remove('adminCredentials');
        router.push("/auth/admin")
    }, [router])
    const { mutateAsync: fetchRefreshToken } = AuthAPI.refreshToken({
        onSuccess: (res) => {
            cookie.setJson("adminCredentials", res.data);
            setCredentials(res.data);
        },
        onError: (error) => {
            // Refresh token has expired...
            signOut();
        }
    })

    useEffect(() => {
        if (credentials) {
            console.log("here")
            if (new Date(credentials.expiresAt * 1000) > new Date()) {
                // handleLoadFetchUSer();
            }
            else {
                fetchRefreshToken(credentials.refreshToken);
            }
        } else {
            signOut()
        }
    }, [])
    const menus = [
        {
            title: "Dashboard",
            href: "/admin",
            Icon: MdDashboard,
        },
        {
            title: "Asset Management",
            href: "/admin/management",
            Icon: IoPersonCircleOutline,
        }, {
            title: "Analytics",
            href: path,
            Icon: CiCircleInfo,
        },
    ]
    const settings = [
        {
            title: "Settings",
            href: "/admin/settings",
            Icon: IoSettingsOutline,
        },
        {
            title: "Log out",
            href: "/admin",
            Icon: CiLogout,
        },
    ]

    const handleClickLogout = () => {
        toast.error("Signed out successfully!")
        cookie.remove('adminCredentials');
        router.push("/auth/admin")
    }


    if (!credentials || new Date(credentials.expiresAt * 1000) < new Date()) {
        return <div></div>
    }


    return (
        <section className="h-screen flex flex-col">
            {
                showLogout && <Modal onTapOutside={() => setShowLogout(false)}>
                    <div onClick={(e) => e.stopPropagation()} className="w-[600px] justify-center items-center h-[300px] rounded-lg bg-zinc-800 flex flex-col">
                        <h3 className="text-white text-[32px]">Confirmation</h3>
                        <p className="text-[20px] text-white">Are you sure you want to log out?</p>
                        <div className="mt-10 flex gap-6 justify-center">
                            <button onClick={handleClickLogout} className="h-[68px] w-[164px] bg-secondary rounded-lg">Yes</button>
                            <button className="h-[68px] text-white w-[164px] ring-border ring-1 rounded-lg" onClick={() => setShowLogout(false)}>No</button>
                        </div>
                    </div>
                </Modal>
            }
            <TopBar />
            <div className="flex-1 flex w-full h-full px-24">
                <section className="w-72 flex flex-col overflow-auto gap-6 pr-3  max-h-[calc(100vh-280px)] py-8 border-r border-secondary h-full justify-between">
                    <ul className="gap-6 flex flex-col">
                        <h4 className="text-secondary text-[24px] font-semibold">MENU</h4>
                        {
                            menus.map(({ title, href, Icon }, index) => index === 2 ? <div className="relative w-full" key={title}>
                                <div className={(isLinkActive(path, href) ? "text-white" : "text-white/40") + " flex gap-4 text-[24px] cursor-pointer items-center text-white"} onClick={() => setShowMore(!showMore)}><Icon />
                                    <span className="text-sm flex-1">{title}</span><IoIosArrowForward className={(showMore && "rotate-90") + ` !text-2xl animate-[translate] text-white/40`} /></div>
                                {
                                    showMore && <ul className="flex flex-col gap-4 mt-4 px-4">
                                        <Link className="text-sm text-white/40 text-[24px]" href="/admin/analytics/site">Site</Link>
                                        <Link className="text-sm text-white/40 text-[24px]" href="/admin/analytics/users">Users</Link>
                                    </ul>
                                }
                            </div> : <Link href={href} key={title} className={(isLinkActive(path, href, index === 0) ? "text-white" : "text-white/40") + " flex gap-4 text-[24px] items-center text-white"}>
                                <Icon />
                                <span className="text-sm flex-1">{title}</span>
                            </Link>)
                        }
                    </ul>
                    <ul className="gap-6 flex flex-col">
                        <h4 className="text-secondary text-[24px] font-semibold">SETTINGS</h4>
                        {
                            settings.map(({ title, href, Icon }, index) => {
                                if (title === "Log out") return <div key={title} className={(isLinkActive(path, href, index === 0) ? "text-white" : "text-white/40") + " cursor-pointer flex gap-4 text-[24px] items-center text-white relative"} onClick={() => setShowLogout(!showLogout)}>
                                    <Icon />
                                    <span className="text-sm">{title}</span>

                                </div>
                                return <Link href={href} key={title} className={(isLinkActive(path, href, index === 0) ? "text-white" : "text-white/40") + " flex gap-4 text-[24px] items-center text-white"}>
                                    <Icon />
                                    <span className="text-sm">{title}</span>
                                </Link>
                            })
                        }
                    </ul>

                </section>
                <section className="flex-1 overflow-auto px-12 py-8 max-h-[calc(100vh-280px)]">
                    {children}
                </section>
            </div>
            <Footer showFull={false} />
        </section>
    );
}
