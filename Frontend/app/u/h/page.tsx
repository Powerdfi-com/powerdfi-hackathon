"use client"
import React, { useCallback, useEffect, useState } from 'react'
import Image from 'next/image';
import { IoWallet } from 'react-icons/io5';
import { FaImage, FaTags } from 'react-icons/fa6';
import Card3 from '@/app/components/cards/card3';
import CustomSlider from '@/app/components/slider';
import Card2 from '@/app/components/cards/card2';
import StatsAPI from '@/app/utils/apis/stats';
import AssetAPI from '@/app/utils/apis/asset';
import Assets from '@/app/components/assets';

const Main = () => {
    const hours = [
        { text: "1hr", value: 1 },
        { text: "24hrs", value: 24 },
        { text: "7days", value: 168 },
        { text: "30days", value: 720 },
    ]
    const steps = [
        {
            title: "Set up your wallet",
            Icon: IoWallet,
        },
        {
            title: "Add your Asset",
            Icon: FaImage
        },
        {
            title: "List them for Sale",
            Icon: FaTags
        }
    ]
    const infos = [
        {
            title: "What are RWAs",
            image: "/card.png"
        },
        {
            title: "What are RWAs",
            image: "/card.png"
        },
        {
            title: "What are RWAs",
            image: "/card.png"
        }
    ]
    const { data: categories, isPending: isFetchingCategories } = AssetAPI.getCategories();
    const { data: chains, isPending: isFetchingChains } = AssetAPI.getChains();
    const defaultFilter = {
        type: "trending",
        blockchain: "",
        range: 24,
    }
    const [filter, setFilter] = useState<{ type: string, blockchain: string; range: number; categoryId?: number }>(defaultFilter);

    const { data: topAssets, isPending: isFetchingTopAssets, isError: isErrorTopAssets, mutate: fetchTopAssets } = StatsAPI.getTopAssets();

    const { data: trendingAssets, isPending: isFetchingTrendingAssets, isError: isErrorTrendingAssets, mutate: fetchTrendingAssets } = StatsAPI.getTrendingAssets();

    useEffect(() => {
        const query = {
            page: 1,
            blockchain: filter.blockchain,
            categoryId: filter.categoryId,
            range: filter.range,
            size: 12,
        };
        if (filter.type === "top") {
            fetchTopAssets(query)
        }
        else {
            fetchTrendingAssets(query)
        }
    }, [filter, fetchTopAssets, fetchTrendingAssets])

    return (
        <main>
            <div className="flex flex-col justify-center items-center py-16 w-full px-48 relative">
                <div className="absolute top-0 bottom-0 right-0 left-0 -z-10">
                    <div className="w-full flex h-full items-center justify-center gap-8">
                        <div className="flex-1 aspect-square gradient rounded-full ring-1 ring-primary/10"></div>
                        <div className="flex-1 aspect-square gradient rounded-full ring-1 ring-primary/10"></div>
                        <div className="flex-1 aspect-square gradient rounded-full ring-1 ring-primary/10"></div>
                    </div>
                </div>
                <Image src="/icon.png" height={150} width={150} alt="logo" />
                <h3 className='text-secondary text-7xl text-center font-semibold max-w-2xl'>The Future <span className='text-white'>of</span> Real World Asset (RWAs)</h3>
            </div>
            <section className="py-12 px-6 sm:px-12 lg:px-24">
                {
                    categories && <ul className='flex gap-4 flex-wrap'>
                        <li><div className={(!filter.categoryId ? 'bg-blue-shade' : '') + (' p-2 cursor-pointer rounded-md text-white text-sm')} onClick={() => setFilter({ ...filter, categoryId: undefined })} >All</div></li>
                        {categories.data.map((category) => <li key={category.id}><div className={(category.id === filter.categoryId ? 'bg-blue-shade' : '') + (' p-2 capitalize cursor-pointer rounded-md text-white text-sm')} onClick={() => setFilter({ ...filter, categoryId: category.id })} >{category.name}</div></li>)}
                    </ul>
                }
                <div className='border-y border-border py-6 mt-8 flex justify-between'>
                    <ul className='flex gap-2 rounded-lg bg-black-shade'>
                        {["Trending", "Top"].map((_tab) => <li key={_tab}><div className={(_tab.toLowerCase() === filter.type ? 'bg-blue-shade rounded-lg' : '') + (' text-white py-2 px-4 text-xs cursor-pointer')} onClick={() => setFilter({ ...filter, type: _tab.toLowerCase() })}>{_tab}</div></li>)}
                    </ul>
                    <div className='flex items-center gap-2'>
                        <ul className='flex gap-2 rounded-lg bg-black-shade'>
                            {hours.map(({ text, value }) => <li key={value}><div className={(value === filter.range ? 'bg-blue-shade rounded-lg' : '') + (' text-white py-2 px-4 text-xs cursor-pointer')} onClick={() => setFilter({ ...filter, range: value })}>{text}</div></li>)}
                        </ul>
                        {
                            chains && <select className='p-2 rounded-lg text-xs text-white capitalize bg-black-shade' defaultValue={filter.blockchain} onChange={(e) => setFilter({ ...filter, blockchain: e.target.value })}>
                                <option value={""}>All Chains</option>
                                {
                                    chains.data.map((chain) => <option key={chain.id} value={chain.id}>{chain.name}</option>)
                                }
                            </select>
                        }
                        <button onClick={() => setFilter(defaultFilter)} className='p-2 text-white text-xs bg-black-shade rounded-lg'>View All</button>
                    </div>
                </div>
                <Assets assets={filter.type === "trending" ? trendingAssets?.data.assets : topAssets?.data.assets} isPending={filter.type === "trending" ? isFetchingTrendingAssets : isFetchingTopAssets} isError={filter.type === "trending" ? isErrorTrendingAssets : isErrorTopAssets} />
            </section>
            <section className="px-24 my-16 ">
                <div className="text-secondary inline font-semibold p-2 rounded-md ring-1 ring-primary w-fit text-3xl">Create RWAs for Sell</div>
                <ul className="grid grid-cols-3 gap-12 mt-10">
                    {steps.map(({ title, Icon }) => <li key={title}>
                        <Card3 Icon={Icon} title={title} />
                    </li>)}
                </ul>
            </section>
            <section className="px-24 my-16 ">
                <div className="text-secondary inline font-semibold p-2 rounded-md ring-1 ring-primary w-fit text-3xl">RWAs 101</div>
                <CustomSlider slidesToShow={3} className='mt-10'>
                    {
                        infos.map(({ title, image }) => <div key={title} className='w-full'><Card2 title={title} image={image} /></div>)
                    }
                </CustomSlider>
            </section>
        </main>
    )
}

export default Main;