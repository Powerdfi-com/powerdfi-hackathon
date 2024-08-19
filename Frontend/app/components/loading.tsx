import Image from 'next/image';
import React, { HTMLAttributes } from 'react'
import { AiOutlineLoading3Quarters } from "react-icons/ai";


const Loading = ({ className }: { className?: HTMLAttributes<HTMLDivElement> | string }) => {
    return (
        <div className={className + ' h-full w-full flex items-center justify-center'}>
            <Image src="/icon.png" alt="icon" height={100} width={100} className='animate animate-pulse' />
        </div>
    )
}

export default Loading;