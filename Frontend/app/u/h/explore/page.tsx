"use client"
import Assets from '@/app/components/assets';
import AssetAPI from '@/app/utils/apis/asset';
import StatsAPI from '@/app/utils/apis/stats';
import { TAsset, TFilter, TStatAssets } from '@/app/utils/types';
import React, { useEffect, useRef, useState } from 'react'
import { FaSpinner } from "react-icons/fa";


const Explore = () => {
    const [topAssets, setTopAssets] = useState<TStatAssets["assets"]>([])
    const [trendingAssets, setTrendingAssets] = useState<TStatAssets["assets"]>([])
    const hours = [
        { text: "1hr", value: 1 },
        { text: "24hrs", value: 24 },
        { text: "7days", value: 7 },
        { text: "30days", value: 30 },
    ]
    const { data: categories, isPending: isFetchingCategories } = AssetAPI.getCategories();
    const { data: chains, isPending: isFetchingChains } = AssetAPI.getChains();
    const defaultFilter = {
        type: "trending",
        blockchain: "",
        range: 24,
    }
    const [totalTop, setTotalTop] = useState(0);
    const [totalTrending, setTotalTrending] = useState(0);
    const [topPage, setTopPage] = useState(1);
    const [trendingPage, setTrendingPage] = useState(1);
    const filterRef = useRef<TFilter>();
    const [filter, setFilter] = useState<TFilter>(defaultFilter);
    const { isPending: isFetchingTopAssets, isError: isErrorTopAssets, mutateAsync: fetchTopAssets, data: topAssetsQuery } = StatsAPI.getTopAssets();

    const { isPending: isFetchingTrendingAssets, isError: isErrorTrendingAssets, mutateAsync: fetchTrendingAssets, data: trendingAssetsQuery } = StatsAPI.getTrendingAssets();

    useEffect(() => {
        const fetchData = async () => {
            const restart = filterRef.current?.categoryId !== filter.categoryId || filterRef.current?.blockchain !== filter.blockchain;
            const query = {
                blockchain: filter.blockchain,
                categoryId: filter.categoryId,
                range: filter.range,
                size: 4,
            };
            if (filter.type === "top") {
                await fetchTopAssets({
                    page: restart ? 1 : topPage,
                    ...query
                }).then((res) => {
                    setTotalTop(res.data.total)
                    setTopAssets((recentAssets) => {
                        if (restart) {
                            return res.data.assets;
                        }
                        return [...recentAssets, ...res.data.assets.filter((asset) => !recentAssets.find((_topAsset) => _topAsset.assetId === asset.assetId))]
                    });
                    if (restart) {
                        setTopPage(1)
                    }
                })
            }
            else {
                await fetchTrendingAssets({
                    page: restart ? 1 : trendingPage,
                    ...query
                }).then((res) => {
                    setTotalTrending(res.data.total)
                    setTrendingAssets((recentAssets) => {
                        if (restart) {
                            return res.data.assets;
                        }
                        else {
                            return [...recentAssets, ...res.data.assets.filter((asset) => !recentAssets.find((_topAsset) => _topAsset.assetId === asset.assetId))]
                        }
                    })
                    if (restart) {
                        setTrendingPage(1)
                    }
                })
            }
            filterRef.current = filter
        }
        fetchData();
    }, [filter, topPage, trendingPage, fetchTopAssets, fetchTrendingAssets])

    return (
        <section className=''>
            <h3 className="text-center text-white text-4xl font-semibold">Explore Assets</h3>
            <section className="py-12 px-24">
                {
                    categories && <ul className='flex gap-4'>
                        <li><div className={(!filter.categoryId ? 'bg-blue-shade' : '') + (' p-2 cursor-pointer rounded-md text-white text-sm')} onClick={() => {
                            setFilter({ ...filter, categoryId: undefined })
                        }} >All</div></li>
                        {categories.data.map((category) => <li key={category.id}><div className={(category.id === filter.categoryId ? 'bg-blue-shade' : '') + (' p-2 cursor-pointer capitalize rounded-md text-white text-sm')} onClick={() => setFilter({ ...filter, categoryId: category.id })} >{category.name}</div></li>)}
                    </ul>
                }
                <div className='border-y border-border py-6 mt-8 flex justify-between'>
                    <ul className='flex gap-2 rounded-lg bg-black-shade'>
                        {["Trending", "Top"].map((_tab) => <li key={_tab}><div className={(_tab.toLowerCase() === filter.type ? 'bg-blue-shade rounded-lg' : '') + (' text-white py-2 px-4 text-xs cursor-pointer')} onClick={() => {
                            setFilter({ ...filter, type: _tab.toLowerCase() });
                        }}>{_tab}</div></li>)}
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
                <Assets assets={filter.type === "trending" ? trendingAssets : topAssets} isPending={filter.type === "trending" ? (isFetchingTrendingAssets && trendingAssets.length === 0) : (isFetchingTopAssets && topAssets.length === 0)} isError={filter.type === "trending" ? isErrorTrendingAssets : isErrorTopAssets} />
            </section>
            <div className='flex w-full justify-center'>
                {
                    ((filter.type === "trending" && trendingAssets.length < totalTrending) || (filter.type === "top" && topAssets.length < totalTop)) ? <button className='cursor-pointer px-4 py-2 rounded-md bg-border text-white min-w-24 text-center flex justify-center items-center' onClick={() => {
                        if (filter.type === "top") {
                            setTopPage((_page) => _page + 1)
                        } else {
                            setTrendingPage((_page) => _page + 1)
                        }
                    }} disabled={(filter.type === "top" && isFetchingTopAssets) || (filter.type === "trending" && isFetchingTrendingAssets)}>{(filter.type === "top" && !isFetchingTopAssets) || (filter.type === "trending" && !isFetchingTrendingAssets) ? <span>Load More</span> : <FaSpinner className='animate animate-spin text-white' />}</button> : <div className='text-white text-center'>Thats all we have for you right now!</div>
                }

            </div>
        </section>
    )
}

export default Explore