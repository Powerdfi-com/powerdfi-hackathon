"use client"
import { Link } from '@/app/components/link';
import AdminAPI from '@/app/utils/apis/admin';
import { TAdmin } from '@/app/utils/types';
import Image from 'next/image';
import { usePathname, useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react'
import { FaSpinner } from 'react-icons/fa6';
import { GoPerson } from "react-icons/go";
import { CiEdit } from "react-icons/ci";



const UsersSettings = () => {
    const { data, mutateAsync, isPending } = AdminAPI.getAdmins();
    const [admins, setAdmins] = useState<TAdmin[]>([]);
    const [totalAdmins, setTotalAdmins] = useState(0)
    const [page, setPage] = useState(1);
    const router = useRouter();
    useEffect(() => {
        const fetchAdmins = async () => {
            await mutateAsync(page).then((res) => {
                setAdmins((prevAdmins) => [...prevAdmins, ...res.data.admins.filter((admin) => !prevAdmins.find((_admin) => admin.id === _admin.id))]);
                setTotalAdmins(res.data.total)
            });
        }
        fetchAdmins()
    }, [page, mutateAsync])
    const path = usePathname()
    return <section>
        <div className='flex justify-between'>
            <div className='text-white flex !text-2xl gap-4'><GoPerson /> <span>Users</span></div>
            <Link href={"/admin/create"} className="flex rounded-full bg-border h-8 items-center relative cursor-pointer">
                <div className="h-8 w-8 rounded-full bg-secondary text-lg flex items-center justify-center">
                    <GoPerson />
                </div>
                <div className="px-4 text-xs text-white">
                    Add Users
                </div>
            </Link>
        </div>
        <ul className='mt-6 grid w-[750px] grid-cols-3 gap-8'>
            {
                admins.map((admin) => <li key={admin.id}>
                    <div className='w-full flex flex-col gap-2'>
                        <div className='w-full h-[184px] relative rounded-lg'>
                            <Image alt={admin.email} fill src="/pool.png" className='rounded-lg' />
                        </div>
                        <div className='w-full flex gap-2'>
                            <div className='flex-1 text-sm text-white'>{admin.email}</div>
                            <Link href={`${path}/edit/${admin.id}`}>
                                <CiEdit className='!text-2xl text-secondary font-semibold' /></Link>
                        </div>
                    </div>
                </li>)
            }
        </ul>
        <div className='flex w-full justify-center mt-6'>
            {
                admins.length < totalAdmins ? <button className='cursor-pointer px-4 py-2 rounded-md bg-border text-white min-w-24 text-center flex justify-center items-center' onClick={() => {
                    setPage((_page) => _page + 1)
                }} disabled={isPending}>{!isPending ? <span>Load More</span> : <FaSpinner className='animate animate-spin text-white' />}</button> : <div className='text-white text-center'>Thats all we have for you right now!</div>
            }

        </div>
    </section>
}

export default UsersSettings