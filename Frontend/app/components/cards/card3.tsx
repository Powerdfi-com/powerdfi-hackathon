import Image from 'next/image'
import React from 'react'
import { IconType } from 'react-icons'

const Card3 = ({ title, Icon }: { title: string, Icon: IconType }) => {
    return (
        <div className='flex flex-col bg-white/10 items-center gap-4 justify-center w-full max-w-sm aspect-video rounded-xl'>
            <div className='h-14 w-14 roundef-full flex items-center justify-center bg-white/40 rounded-full'>
                <Icon className='text-secondary !text-xl font-semibold' />
            </div>
            <h4 className='leading-relaxed text-white font-semibold'>{title}</h4>
        </div>
    )
}

export default Card3