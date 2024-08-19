"use client"
import { Link } from '@/app/components/link';
import React from 'react'
import { GoPerson } from "react-icons/go";
import { PiNotificationFill } from "react-icons/pi"; import { FaShieldAlt } from "react-icons/fa";




const Settings = () => {
    return <section>
        <h3 className='text-white text-3xl'>Setting</h3>
        <ul className='flex gap-8 flex-wrap mt-6'>
            <Link href={"/admin/settings/general"} className='w-[150px] h-[160px] ring-1 ring-border flex flex-col items-center justify-center gap-4 rounded-sm'>
                <FaShieldAlt className='text-white !text-3xl' />
                <div className='text-white'>General</div>
            </Link>
            <Link href={"/admin/settings/users"} className='w-[150px] h-[160px] ring-1 ring-border flex flex-col items-center justify-center gap-4 rounded-sm'>
                <GoPerson className='text-white !text-3xl' />
                <div className='text-white'>Users</div>
            </Link>
            <Link href={"/admin/settings/notifications"} className='w-[150px] h-[160px] ring-1 ring-border flex flex-col items-center justify-center gap-4 rounded-sm'>
                <PiNotificationFill className='text-white !text-3xl' />
                <div className='text-white'>Notification</div>
            </Link>
        </ul>
    </section>
}

export default Settings