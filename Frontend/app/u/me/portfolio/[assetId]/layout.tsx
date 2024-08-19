"use client"
import Error from '@/app/components/error'
import { Link } from '@/app/components/link'
import Loading from '@/app/components/loading'
import { AssetContext } from '@/app/context/context'
import AssetAPI from '@/app/utils/apis/asset'
import { isLinkActive } from '@/app/utils/func'
import Image from 'next/image'
import { usePathname } from 'next/navigation'
import React, { ReactNode } from 'react'

const PortfolioLayout = ({ children, params }: { children: ReactNode, params: { assetId: string } }) => {
    const links = [
        {
            text: "About",
            href: `/u/me/portfolio/${params.assetId}`
        }, {
            text: "Property Type",
            href: `/u/me/portfolio/${params.assetId}/property`
        }, {
            text: "Transaction History",
            href: `/u/me/portfolio/${params.assetId}/history`
        },
    ]
    const path = usePathname();
    const { isPending, data } = AssetAPI.getAssetById(params.assetId);
    return <section>
        <h3 className='text-2xl text-white leading-relaxed'>Manage your Asset Token</h3>
        {
            isPending ? <Loading /> : (data ? <>
                <div className='mt-8 flex gap-6'>
                    <div className="flex-1 h-64">
                        <div className='h-64 w-full relative'>
                            <Image src="/item.png" alt="alt" fill={true} className='object-cover rounded-xl' />
                        </div>
                    </div>
                    <div className='flex-[2] ring-1 ring-primary/40 rounded-2xl p-6 gradient'>
                        <div className='flex justify-between'>
                            <div>
                                <h3 className='text-white text-xl leading-relaxed'>Elvon Eve</h3>
                                <p className='text-secondary text-sm'>Elvon Ave</p>
                            </div>
                        </div>
                        <div className="flex gap-4 mt-4">
                            <div className='bg-blue-shade rounded-xl p-2 flex items-center justify-center flex-col flex-1 gap-1'>
                                <div className='text-white/40 text-sm '>Token Own</div>
                                <div className='flex gap-1 items-center'><span className="text-white text-2xl">200PDFi</span><span className='text-white/40 text-sm'>$2,470. 89</span></div>
                            </div>
                            <div className='bg-blue-shade rounded-xl p-2 flex items-center justify-center flex-col flex-1 gap-1'>
                                <div className='text-white/40 text-sm '>Token Own</div>
                                <div className='flex gap-1 items-center'><span className="text-white text-2xl">200PDFi</span><span className='text-white/40 text-sm'>$2,470. 89</span></div>
                            </div>
                        </div>
                        <div className="flex gap-2 w-full mt-6">
                            <button className="text-sm bg-secondary h-8 rounded-md flex-1">Invest</button>
                            <button className="text-sm bg-transparent ring-1 ring-primary/40 text-secondary h-8 rounded-md flex-1">Stake your Investment</button>
                        </div>
                    </div>
                </div>
                <section className='mt-8'>
                    <ul className="flex gap-10">
                        {links.map(({ text, href }) => <li key={text}>
                            <Link href={href} className={(isLinkActive(href, path) ? "text-secondary" : "text-white") + ' p-1.5 text-[24px]'}>{text}</Link>
                        </li>)}
                    </ul>
                    <section className='my-6'>
                        <AssetContext.Provider value={data?.data}>
                            {children}
                        </AssetContext.Provider>
                    </section>
                </section>
            </> : <Error />)
        }
    </section>
}

export default PortfolioLayout