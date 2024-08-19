"use client"
import Error from '@/app/components/error';
import Loading from '@/app/components/loading';
import AssetAPI from '@/app/utils/apis/asset'
import ListingAPI from '@/app/utils/apis/listing';
import { dateToISO } from '@/app/utils/func';
import Image from 'next/image'
import { useRouter } from 'next/navigation';
import React, { useState } from 'react'
import { toast } from 'react-toastify';

const ListAsset = ({ params }: { params: { assetId: string, chain: string } }) => {
    const router = useRouter()
    const { data: asset, isPending } = AssetAPI.getAssetById(params.assetId);
    const { mutateAsync: createLisitng } = ListingAPI.createListing({
        onSuccess: () => {
            router.push("/u/me/listing")
        }
    });
    const { data: currencies, isPending: isFetchingCurrencies } = ListingAPI.getChainTokens(params.chain);
    const [data, setData] = useState<{
        currency: string[];
        price: string;
        quantity: string;
        startAt: string;
        endAt: string;
        min_investment_amount: string;
        max_investment_amount: string;
        max_raise_amount: string;
        min_raise_amount: string;
        type: string
    }>({
        currency: [],
        endAt: "",
        max_investment_amount: "",
        max_raise_amount: "",
        min_investment_amount: "",
        min_raise_amount: "",
        price: "",
        quantity: "",
        startAt: "",
        type: "auction"
    });

    const handleClickCreateLisitng = async (e: any) => {
        e.preventDefault();
        await toast.promise(createLisitng({ ...data, assetId: params.assetId, price: parseInt(data.price), quantity: parseInt(data.quantity), min_investment_amount: parseInt(data.min_investment_amount), max_investment_amount: parseInt(data.max_investment_amount), min_raise_amount: parseInt(data.min_raise_amount), max_raise_amount: parseInt(data.max_raise_amount), startAt: dateToISO(data.startAt), endAt: dateToISO(data.endAt) }), {
            error: "Something went wrong!",
            pending: "Listing asset... Please wait!",
            success: "Asset Listed Successfully!"
        },)
    }

    return <section>
        <h3 className='text-2xl text-white leading-relaxed'>Token Listing Information</h3>
        {
            isPending ? <Loading /> : (asset ? <div className='mt-8 flex gap-6'>
                <div className="flex-1">
                    <div className='aspect-video w-full relative'>
                        <Image src={asset.data.urls[0]} alt={asset.data.name} fill={true} className='object-cover rounded-xl' />
                    </div>
                </div>
                <form className='flex-1 flex flex-col gap-2'>
                    <label className='flex flex-col gap-2'>
                        <span className='text-sm text-white'>Asset Listing Name</span>
                        <input className='text-white max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' disabled value={asset.data.name} />
                    </label>
                    <label className='flex flex-col gap-2'>
                        <span className='text-sm text-white'>Amount of token to be mint</span>
                        <input type='number' className='text-white max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' placeholder='Elton Ave' value={data.quantity} onChange={(e) => setData({ ...data, quantity: e.target.value })} />
                    </label>
                    <div className='flex gap-4 w-full max-w-sm'>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>Starting Date</span>
                            <input type='date' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.startAt} onChange={(e) => setData({ ...data, startAt: e.target.value })} />
                        </label>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>End Date</span>
                            <input type='date' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.endAt} onChange={(e) => setData({ ...data, endAt: e.target.value })} />
                        </label>
                    </div>
                    <div className='flex gap-4 w-full max-w-sm'>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>Min amount to raise</span>
                            <input type='number' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.min_raise_amount} onChange={(e) => setData({ ...data, min_raise_amount: e.target.value })} />
                        </label>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>Max amount to raise</span>
                            <input type='number' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.max_raise_amount} onChange={(e) => setData({ ...data, max_raise_amount: e.target.value })} />
                        </label>
                    </div>
                    <label className='flex flex-col gap-2 mt-3'>
                        <span className='text-sm text-white'>ERC-20 Tokens Accepted</span>
                        <select multiple className='max-w-sm bg-transparent border-text-grey border rounded-md text-white py-2 text-sm px-5 ' value={data.currency} onChange={(e) => setData({ ...data, currency: Array.from(e.target.selectedOptions, option => option.value) })} >
                            <option>Select Currency</option>
                            {
                                currencies && currencies.data.map((currency) => <option key={currency.id} value={currency.id}>{currency.name}</option>)
                            }
                        </select>
                    </label>
                    <label className='flex flex-col gap-2'>
                        <span className='text-sm text-white'>Token Price in USD</span>
                        <input type='number' className='max-w-sm outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' placeholder='Elton Ave' value={data.price} onChange={(e) => setData({ ...data, price: e.target.value })} />
                    </label>
                    <div className='flex gap-4 w-full max-w-sm'>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>Min amount in USD</span>
                            <input type='number' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.min_investment_amount} onChange={(e) => setData({ ...data, min_investment_amount: e.target.value })} />
                        </label>
                        <label className='flex-1 flex flex-col gap-2'>
                            <span className='text-sm text-white'>Max amount in USD</span>
                            <input type='number' className='w-full outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5 text-white' value={data.max_investment_amount} onChange={(e) => setData({ ...data, max_investment_amount: e.target.value })} />
                        </label>
                    </div>
                    <div className="flex gap-2 w-full mt-6 mb-6">
                        <button className="text-sm bg-transparent ring-1 ring-primary/40 text-secondary h-8 rounded-md flex-1">Cancel</button>
                        <button className="text-sm bg-secondary h-8 rounded-md flex-1" onClick={handleClickCreateLisitng}>Complete Listing</button>
                    </div>
                </form>
            </div> : <Error />)
        }
    </section>
}

export default ListAsset