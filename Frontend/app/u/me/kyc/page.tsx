"use client"
import Modal from '@/app/components/modal';
import UserAPI from '@/app/utils/apis/user'
import React, { useContext, useEffect, useState } from 'react'
import { toast } from 'react-toastify';
import { MdOutlineVerified } from "react-icons/md";
import { UserContext } from '@/app/context/context';
import { Link } from '@/app/components/link'
import { LuScanFace } from "react-icons/lu";
import { IoDocumentOutline } from "react-icons/io5"; import { HiArrowLongRight, HiArrowLongLeft, HiArrowLongUp, HiArrowLongDown } from "react-icons/hi2";
import { CiMail } from "react-icons/ci";







const KYC = () => {
    const { user } = useContext(UserContext);
    const [link, setLink] = useState("");
    const { mutateAsync: getKycLink, isPending } = UserAPI.getKycLink({
        onSuccess: (res) => {
            console.log(res.data.link)
            setLink(res.data.link);
        }
    });
    const handleClickVerify = async () => {
        await toast.promise(getKycLink, {
            pending: "Please wait...",
            success: "Please conitinue verification on the new modal!",
            error: "Something went wrong, please try again later!"
        });
    }
    const callback = (e: any) => {
        console.log(e)
    }
    useEffect(() => {
        window.addEventListener('message', callback);
        return () => window.removeEventListener('message', callback);
    }, []);

    return user.isVerified ? <IsVerified /> : <section>
        {
            link && <Modal onTapOutside={() => setLink("")}>
                <iframe src={link} className="w-full h-full bg-white"></iframe>
            </Modal>
        }
        <h3 className='text-4xl text-white leading-relaxed text-center font-semibold'>Verify your Identity</h3>
        <div className='flex flex-col  mt-8 gap-4'>
            <div className='flex gap-2 items-center'>
                <div className='w-[263px] h-[160px] rounded-xl flex ring-1 ring-border flex-col gap-4 items-center justify-center'>
                    <LuScanFace className='text-secondary font-semibold text-6xl' />
                    <p className='text-md text-white'>Facial Recognition</p>
                </div>
                <HiArrowLongRight className='!text-7xl text-secondary font-semibold' />
                <div className='w-[263px] h-[160px] rounded-xl flex ring-1 ring-border flex-col gap-4 items-center justify-center'>
                    <CiMail className='text-secondary font-semibold text-6xl' />
                    <p className='text-md text-white'>Email Verification</p>
                </div>
                <HiArrowLongRight className='!text-7xl text-secondary font-semibold' />
                <div className='w-[263px] h-[160px] rounded-xl flex ring-1 ring-border flex-col gap-4 items-center justify-center'>
                    <LuScanFace className='text-secondary font-semibold text-6xl' />
                    <p className='text-md text-white'>Verification Successful </p>
                </div>
            </div>
            <div className='flex justify-between'>
                <div className='w-[263px] flex items-center justify-center'>
                    <HiArrowLongUp className='!text-7xl text-secondary font-semibold' />
                </div>
                <div className='w-[263px] flex items-center justify-center'>
                    <HiArrowLongDown className='!text-7xl text-secondary font-semibold' />
                </div>
            </div>
            <div className='flex gap-2 items-center'>
                <div className='w-[263px] h-[160px] rounded-xl flex ring-1 ring-border flex-col gap-4 items-center justify-center'>
                    <IoDocumentOutline className='text-secondary font-semibold text-6xl' />
                    <p className='text-md text-white'>ID Verification</p>
                </div>
                <HiArrowLongLeft className='!text-7xl text-secondary font-semibold' />
                <div className='w-[263px] h-[160px] rounded-xl flex ring-1 ring-border flex-col gap-4 items-center justify-center'>
                    <IoDocumentOutline className='text-secondary font-semibold text-6xl' />
                    <p className='text-md text-white'>AML</p>
                </div>
                <HiArrowLongLeft className='!text-7xl text-secondary font-semibold' />
                <button className="w-[276px] self-center h-[68px] rounded-lg bg-secondary font-semibold text-lg" onClick={handleClickVerify}>Verify your Account</button>
            </div>
        </div>
    </section>
}

export default KYC


const IsVerified = () => {
    return <section className="w-full flex h-full items-center justify-center">
        <article className="flex flex-col gap-4 max-w-sm items-center justify-center">
            <MdOutlineVerified className='!text-9xl text-secondary' />
            <h4 className='text-white text-md font-semibold leading-relaxed'>Account Verified</h4>
            <p className='text-center text-white/40 text-sm'>Your KYC details have been verified successfully. Click the button to proceed with the creation of Asset Token.</p>
            <Link href="/u/me/create/chains" className='w-64 h-10 rounded-md bg-secondary flex items-center justify-center'>Create Token</Link>
        </article>
    </section>
}