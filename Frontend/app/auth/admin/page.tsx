"use client"
import Image from 'next/image'
import React, { useState } from 'react'
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import AuthAPI from '@/app/utils/apis/auth';
import { cookie } from '@/app/utils/cookie';

const AdminSignIn = () => {
    const router = useRouter();
    const [data, setData] = useState({
        password: "",
        email: "",
    });
    const { mutateAsync: signInAdmin, isPending: isSigningIn } = AuthAPI.adminSignIn({
        onSuccess: (res) => {
            console.log(res.data);
            cookie.setJson('adminCredentials', res.data.tokens);
            router.push("/admin/");
        },
    });
    const handleClickAdminSignIn = async () => {
        await toast.promise(signInAdmin(data), {
            error: "Something went wrong, please try again!",
            pending: "Signing in.. please wait!",
            success: "Signed in, Please continue!"
        })
    }
    return <div className='m-24 px-8 sm:px-12 md:px-24 items-center flex rounded-xl ring-1 ring-primary/40 bg-black-shade overflow-auto' onClick={(e) => e.stopPropagation()}>
        <div className='flex-[2]'>
            <Image src={"/icon.png"} alt="icon" height={100} width={100} />
            <h3 className='text-white font-semibold text-2xl'>Sign In</h3>
            <p className='text-white text-sm max-w-lg leading-relaxed mt-2'>
                Enter your email address and password to enable you to verify or reject asset.
            </p>
            <form>
                <h4 className='text-white mt-8 mb-4'>Personal Information</h4>
                <label className='flex flex-col gap-1 mt-3'>
                    <span className='text-sm text-white'>Email</span>
                    <input type='email' value={data.email} onChange={(e) => setData({ ...data, email: e.target.value })} className='max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-3 px-5 text-white text-sm' placeholder='First Name' />
                </label>
                <label className='flex flex-col gap-1 mt-3'>
                    <span className='text-sm text-white'>Password</span>
                    <input value={data.password} onChange={(e) => setData({ ...data, password: e.target.value })} className='max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-3 px-5 text-white text-sm' type='password' placeholder='PowerDfi@gmail.com' />
                </label>

            </form>

            <button disabled={isSigningIn} onClick={handleClickAdminSignIn} className='h-12 w-[200px] bg-secondary text-sm rounded-md mt-8 mb-24'>Login</button>
        </div>
        <div className='flex h-full items-center justify-center flex-1'>
            <p className='text-secondary text-6xl leading-relaxed max-w-md font-semibold '>The Future <span className='text-white'>of</span> Real World Asset (RWA)</p>
        </div>
    </div>
}

export default AdminSignIn