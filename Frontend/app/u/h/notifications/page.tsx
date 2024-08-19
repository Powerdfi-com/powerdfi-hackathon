"use client"
import Loading from '@/app/components/loading';
import NotificationsAPI from '@/app/utils/apis/notifications';
import { TNotifications } from '@/app/utils/types';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react'
import { FaSpinner } from 'react-icons/fa6';
import { IoMdCheckmark, IoMdClose } from 'react-icons/io'
import { toast } from 'react-toastify';

const Notifications = () => {
    const router = useRouter();
    const [page, setPage] = useState(1);
    const [notifications, setNotifications] = useState<TNotifications["notifications"]>([])
    const { mutateAsync: fetchNotifications, isPending: isPendingNotifications } = NotificationsAPI.getNotifications();
    const { mutateAsync, isPending } = NotificationsAPI.markAllAsRead({
        onSuccess: () => {
            router.refresh()
        }
    });
    const [totalNotifications, setTotalNotifications] = useState(0)

    const handleClickMarkAllAsRead = () => [
        toast.promise(mutateAsync(), {
            pending: "Working... Please wait!",
            error: "Something went wrong, please try again later!",
            success: "All notifications marked as read!"
        })
    ]

    useEffect(() => {
        const fetchData = async () => {
            await fetchNotifications({ page }).then((res) => {
                setNotifications((prevNotifcations) => {
                    return [...prevNotifcations, ...res.data.notifications.filter((_notification) => !prevNotifcations.find((notification) => notification.id === _notification.id))]
                });
                setTotalNotifications(res.data.total)
            });
        }
        fetchData()
    }, [page, fetchNotifications])

    return <section className='flex justify-center w-full'>
        <section className='max-w-[800px] w-full px-6'>
            <div className='flex w-full gap-6'>
                <div className='flex flex-col gap-4 justify-between flex-1'>
                    <h3 className='text-white text-4xl font-semibold'>Notification</h3>
                    <p className='text-white/60 text-md'>You have <span className='text-secondary'>{notifications.filter((notification) => !notification.viewed).length}</span> notifications to go through</p>
                </div>
                <button className="bg-secondary disabled:bg-border h-[45px] px-4 rounded-md text-xl relative" disabled={isPending || !notifications.find((notification) => !notification.viewed)} onClick={handleClickMarkAllAsRead}>Mark all as read</button>
            </div>
            {
                notifications ? <>
                    <section className='mt-6'>
                        <h4 className='text-2xl text-white/60'>Today</h4>
                        {
                            !notifications ? <Loading /> : <ul className="flex flex-col gap-4 mt-6">
                                {
                                    notifications.map((notification) => <li key={notification.id}>
                                        <div className={`${!notification.viewed && "bg-border"} w-full h-[96px] ring-border ring-1 rounded-md flex items-center px-4 gap-4`}>
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
                    </section>
                    <div className='flex w-full justify-center mt-6'>
                        {
                            notifications.length < totalNotifications ? <button className='cursor-pointer px-4 py-2 rounded-md bg-border text-white min-w-24 text-center flex justify-center items-center' onClick={() => {
                                setPage((_page) => _page + 1)
                            }} disabled={isPendingNotifications}>{!isPendingNotifications ? <span>Load More</span> : <FaSpinner className='animate animate-spin text-white' />}</button> : <div className='text-white text-center'>Thats all we have for you right now!</div>
                        }

                    </div>

                </> : <Loading />
            }
        </section>
    </section>
}

export default Notifications

const AccpetedNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-green-800 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-green-400 rounded-full flex items-center justify-center">
            <IoMdCheckmark className="text-black/40 text-3xl" />
        </div>
    </div>
}

const RejectedNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-red-800 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-red-400 rounded-full flex items-center justify-center">
            <IoMdClose className="text-black/40 text-3xl" />
        </div>
    </div>
}

const SoldNotification = () => {
    return <div className="h-[40px] w-[40px] rounded-full bg-cyan-800 flex items-center justify-center">
        <div className="h-[27px] w-[27px] bg-cyan-400 rounded-full flex items-center justify-center">
            <IoMdClose className="text-black/40 text-3xl" />
        </div>
    </div>
}
