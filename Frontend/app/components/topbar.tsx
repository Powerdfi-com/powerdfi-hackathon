"use client"
import Image from "next/image"
import { CgNotifications } from "react-icons/cg";
import { LuWallet2 } from "react-icons/lu";
import { IoMdClose } from "react-icons/io";
import { IoBriefcase } from "react-icons/io5";
import { GrTransaction } from "react-icons/gr";
import { IoGrid } from "react-icons/io5";
import { IoSettings } from "react-icons/io5";
import { Dispatch, SetStateAction, useContext, useEffect, useState } from 'react'
import { HiMenuAlt3 } from "react-icons/hi";
import { IoCopy } from "react-icons/io5";
import { toast } from "react-toastify";
import { Link } from "./link";
import { MediaQueryContext, UserContext } from "../context/context";
import UserAPI from "../utils/apis/user";
import { IoMdCheckmark } from "react-icons/io";
import NotificationsAPI from "../utils/apis/notifications";
import Loading from "./loading";





const TopBar = ({ home, showWallet, showProfile, setShowProfile, showNotifications, setShowWallet, setShowNotifications }: { home: string, showWallet: boolean, showNotifications: boolean, showProfile: boolean, setShowProfile: Dispatch<SetStateAction<boolean>>, setShowWallet: Dispatch<SetStateAction<boolean>>, setShowNotifications: Dispatch<SetStateAction<boolean>> }) => {
    const navLinks = [
        {
            text: "Explore",
            href: "/u/h/explore"
        }, {
            text: "Create",
            href: "/u/me/create"
        }, {
            text: "Stats",
            href: "/u/h/stats"
        },
    ];
    const notificationTabs = [
        {
            text: 'Portfolio',
            href: "/u/me/portfolio",
            Icon: IoBriefcase,
        },
        {
            text: 'Transactions',
            href: "/u/me/transactions",
            Icon: GrTransaction,
        }, {
            text: 'Explore',
            href: "/u/h/explore",
            Icon: IoGrid,
        }, {
            text: 'Settings',
            href: "/u/me/settings",
            Icon: IoSettings,
        }
    ]
    const { signOut, copyAddress } = useContext(UserContext)

    const { user, address } = useContext(UserContext);
    const { data } = UserAPI.getWalletDetails();
    const { data: notifications, isPending, mutate: fetchNotifications } = NotificationsAPI.getNotifications();

    useEffect(() => {
        fetchNotifications({ size: 5, page: 1 })
    }, [fetchNotifications])

    const isMobile = useContext(MediaQueryContext);
    const [showNav, setShowNav] = useState(false)


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
        <section className="flex items-center px-8 sm:px-12 md:px-24 py-8 gap-12">
            <Link href={"/i"}><Image src="/logo.png" alt="logo" height={50} width={50} /></Link>
            <nav className="flex-1">
                <ul className="flex gap-8">
                    {navLinks.map(({ text, href }) => <li key={text}><Link href={href} className="text-white text-md font-semibold">{text}</Link></li>)}
                </ul>
            </nav>
            <div className="flex gap-2 items-center">
                <div className="flex rounded-full bg-black-shade h-8 items-center relative cursor-pointer" onClick={(e) => { setShowNotifications(!showNotifications); setShowWallet(false); setShowProfile(false); e.stopPropagation() }}>
                    <div className="h-8 w-8 rounded-full bg-secondary text-lg flex items-center justify-center">
                        <CgNotifications />
                    </div>
                    <div className="px-4 text-xs text-white">
                        Notifications
                    </div>
                    {
                        showNotifications && <div className="cursor-default absolute z-10 w-[500px] ring-1 p-8 bg-black-shade rounded-lg ring-primary/40 top-12 right-0" onClick={(e) => e.stopPropagation()}>
                            <h3 className="text-white text-3xl">Notification</h3>
                            {
                                isPending ? <Loading /> :
                                    <ul className="flex flex-col gap-4 mt-6">
                                        {
                                            notifications?.data.notifications.map((notification) => <li key={notification.id}>
                                                <div className="w-full h-[52px] ring-border ring-1 rounded-md flex items-center px-4 gap-4">
                                                    {
                                                        notification.type === "approve" ? <AccpetedNotification /> : (notification.type === "reject" ? <RejectedNotification /> : <SoldNotification />)
                                                    }
                                                    <span className="text-white text-[17px]">{notification.type === "approve" ? `New Asset "${notification.data.assetName}" has been verified!` : (
                                                        notification.type === "reject" ? `New Asset "${notification.data.assetName}" has been rejected!` : `New Asset "${notification.data.assetName}" has been sold!`
                                                    )}</span>
                                                </div>
                                            </li>)
                                        }
                                    </ul>
                            }
                            <div className="mt-12 flex justify-center">
                                <Link onClick={(e) => { setShowWallet(false); setShowNotifications(false); setShowProfile(false); e.stopPropagation() }} href={"/u/h/notifications"} className="bg-secondary px-8 py-3 rounded-md text-xl relative">View all Notifications</Link>
                            </div>
                        </div>
                    }
                </div>
                <div className="flex rounded-full cursor-pointer bg-black-shade h-8 items-center relative" onClick={(e) => { setShowWallet(!showWallet); setShowNotifications(false); setShowProfile(false); e.stopPropagation() }}>
                    {
                        showWallet && <div className="cursor-default absolute z-10 w-96 ring-1 p-8 bg-black-shade rounded-lg ring-primary/40 top-12 right-0" onClick={(e) => e.stopPropagation()}>
                            <div className="flex gap-6 flex-col">

                                <div className="text-md font-semibold text-white flex-1">
                                    Your Wallet
                                </div>
                                <div className="flex w-full gap-4 items-center">
                                    <div className="h-14 w-14 relative aspect-square ring-1 ring-primary/40 rounded-md flex items-center justify-center">
                                        <div className="relative h-10 w-10">
                                            <Image src="/icon.png" alt="powerdfi" fill={true} className="object-cover" />
                                        </div>
                                    </div>
                                    <div className="text-sm text-white overflow-ellipsis overflow-hidden flex-1">{address}</div>
                                    <button className="text-white !font-semibold !text-lg" onClick={copyAddress}><IoCopy /></button>
                                </div>
                                {
                                    data && <div className="w-full mt-6 rounded-md ring-1 ring-primary/40 py-4 flex flex-col gap-2 items-center justify-center">
                                        <div className="font-light text-white/40 text-xs">Available Balance</div>
                                        <div className="text-2xl text-white">$ {data.data.balance}</div>
                                    </div>
                                }
                                <div className="flex gap-2 w-full">
                                    <button className="text-sm h-10 rounded-md flex-1 bg-secondary">Add Funds</button>
                                    <button className="text-sm text-white h-10 rounded-md flex-1 ring-1 ring-primary/40" onClick={signOut}>Disconnect Wallet</button>
                                </div>
                            </div>
                        </div>
                    }
                    <div className="h-8 w-8 rounded-full bg-secondary text-lg flex items-center justify-center">
                        <LuWallet2 />
                    </div>
                    <div className="px-4 text-xs text-white">
                        Wallet
                    </div>
                </div>
                <div className="h-10 w-10 relative bg-black-shade rounded-full" onClick={(e) => { setShowProfile(!showProfile); setShowWallet(false); setShowNotifications(false); e.stopPropagation() }}>
                    <Image src={user.avatar || "/avatar.png"} objectFit="cover" alt="user" fill={true} className="cursor-pointer rounded-full" />
                    {
                        showProfile && <div className="cursor-default absolute z-10 w-72 ring-1 p-8 bg-black-shade rounded-lg ring-primary/40 top-12 right-0" onClick={(e) => e.stopPropagation()}>
                            <div className="flex gap-4 items-center">
                                <div className="w-10 h-10 relative">
                                    <Image src={user.avatar || "/avatar.png"} alt="user" fill={true} className="object-cover rounded-full" />
                                </div>
                                <div className="text-sm text-white flex-1 overflow-hidden overflow-ellipsis">
                                    {address}
                                </div>
                                <button onClick={() => setShowProfile(false)}><IoMdClose className="text-primary !text-xl" /></button>
                            </div>
                            <ul className="flex flex-col gap-4 my-6">
                                {notificationTabs.map(({ text, href, Icon }) => <li key={href}>
                                    <Link href={href} className="flex items-center gap-4"><Icon className="!text-xl text-secondary" /><span className="text-sm text-white font-light">{text}</span></Link>
                                </li>)}
                            </ul>
                        </div>
                    }
                </div>
            </div>
        </section >
    )
}

export default TopBar

const AccpetedNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-green-200 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-green-400 rounded-full flex items-center justify-center">
            <IoMdCheckmark className="text-black/40 text-3xl" />
        </div>
    </div>
}

const RejectedNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-red-200 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-red-400 rounded-full flex items-center justify-center">
            <IoMdClose className="text-black/40 text-3xl" />
        </div>
    </div>
}

const SoldNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-cyan-200 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-cyan-400 rounded-full flex items-center justify-center">
            <IoMdClose className="text-black/40 text-3xl" />
        </div>
    </div>
}
