import Image from 'next/image'
import React from 'react'

const Verified = () => {
    return <div className='relative h-10 w-10'>
        <Image src="/verified.png" alt="verified" fill className='object-contain' />
    </div>
}

export default Verified