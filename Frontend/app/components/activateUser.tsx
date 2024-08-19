"use client"
import Image from 'next/image'
import React, { useState } from 'react'
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import UserAPI from '@/app/utils/apis/user';
import { cookie } from '@/app/utils/cookie';
import { Link } from '@/app/components/link';


const ActivateUser = () => {
    const router = useRouter();
    const [data, setData] = useState({
        username: "",
        email: "",
    });
    const { mutateAsync: activateUser, isPending: isActivatingUser } = UserAPI.activateUser({
        onSuccess: (res) => {
            cookie.setJson('credentials', res.data.tokens);
            router.push("/u/h");
        },
    });
    const handleClickActivateUser = async () => {
        await toast.promise(activateUser(data), {
            error: "Something went wrong, please try again!",
            pending: "Activating user... please wait!",
            success: "Your account has been activated!"
        })
    }
    return <div className='w-full p-24 flex h-full rounded-xl ring-1 ring-primary/40 bg-black-shade overflow-auto' onClick={(e) => e.stopPropagation()}>
        <div className='flex-[2]'>
            <Image src={"/icon.png"} alt="icon" height={100} width={100} />
            <h3 className='text-white font-semibold text-2xl'>Wallet Connected</h3>
            <p className='text-white text-sm max-w-lg leading-relaxed mt-2'>
                Enter your username and email address to enable you create and manage your security token.
            </p>
            <form>
                <h4 className='text-white mt-8 mb-4'>Personal Information</h4>
                <label className='flex flex-col gap-1 mt-3'>
                    <span className='text-sm text-white'>Username</span>
                    <input value={data.username} onChange={(e) => setData({ ...data, username: e.target.value })} className='max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-3 px-5 text-white text-sm' placeholder='First Name' />
                </label>
                <label className='flex flex-col gap-1 mt-3'>
                    <span className='text-sm text-white'>Email</span>
                    <input value={data.email} onChange={(e) => setData({ ...data, email: e.target.value })} className='max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-3 px-5 text-white text-sm' type='email' placeholder='PowerDfi@gmail.com' />
                </label>

            </form>
            <p className='mt-6 text-white/40 text-sm max-w-sm'>
                By using PowerDfi, you agree to our <Link href={"/"}>Terms of Service</Link> and our <Link href="/">Privacy Policy</Link>.
            </p>
            <button disabled={isActivatingUser} onClick={handleClickActivateUser} className='h-12 w-full max-w-xs bg-secondary text-sm rounded-md mt-8 mb-24'>Create Account</button>
        </div>
        <div className='flex h-full items-center justify-center flex-1'>
            <p className='text-secondary text-6xl leading-relaxed max-w-md font-semibold '>The Future <span className='text-white'>of</span> Real World Asset (RWA)</p>
        </div>
    </div>
}

export default ActivateUser