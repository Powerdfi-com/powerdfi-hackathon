"use client"
import React, { ReactNode, useContext, useState } from 'react'
import { MdFavorite } from "react-icons/md";
import { IoEyeSharp } from "react-icons/io5";
import { IoMdShare } from "react-icons/io";
import Card1 from '@/app/components/cards/card1';
import { TAsset } from '@/app/utils/types';
import Image from 'next/image'; import { FaDiscord } from "react-icons/fa";
import { FaXTwitter } from "react-icons/fa6";
import { FaExternalLinkAlt } from "react-icons/fa";
import AssetAPI from '@/app/utils/apis/asset';
import Loading from '@/app/components/loading';
import Error from '@/app/components/error';
import { Link } from '@/app/components/link';
import { usePathname } from 'next/navigation';
import { isLinkActive } from '@/app/utils/func';
import { toast } from 'react-toastify';
import { IoMdRemove } from "react-icons/io";
import { IoMdAdd } from "react-icons/io";
import OrderAPI from '@/app/utils/apis/order';
import { QueryClient } from '@tanstack/react-query';

import { BsAward } from "react-icons/bs";
import { UserContext } from '@/app/context/context';
import Modal from '@/app/components/modal';




const Asset = ({ children, params }: { children: ReactNode, params: { assetId: string } }) => {
    const { isPending, data } = AssetAPI.getAssetById(params.assetId);
    const navs = [
        {
            text: "Details",
            href: `/u/h/asset/${params.assetId}`
        },
        {
            text: "Order Book",
            href: `/u/h/asset/${params.assetId}/order-book`
        },
        {
            text: "Documents",
            href: `/u/h/asset/${params.assetId}/documents`
        }
    ]
    const { user } = useContext(UserContext);
    const path = usePathname();
    return isPending ? <Loading /> : (data ? <section className="flex flex-col gap-16  px-6 sm:px-12 lg:px-24">
        <section className="flex justify-between relative">
            <article className="flex flex-col gap-2">
                <h3 className="text-white text-[40px]">{data.data.name}  </h3>
                <p className='text-sm text-white/80'>{data.data.name}</p>
            </article>
            <div className="flex gap-4 items-center">
                <div className="flex items-center gap-2 text-sm text-text-grey2">
                    <MdFavorite className="!text-lg" />
                    <span className="font-light">{data.data.favourites} Favorites</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-text-grey2">
                    <IoEyeSharp className="!text-lg" />
                    <span className="font-light">{data.data.views} Views</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-text-grey2">
                    <IoMdShare className="!text-lg" />
                    <span className="font-light">Share</span>
                </div>
            </div>
        </section>
        <section className='bg-black-shade p-4 rounded-lg flex flex-col gap-4 h-[760px] relative'>
            <div className="flex-1 flex w-full gap-8">
                <div className='flex-1 w-full relative'>
                    <Image src={data.data.urls[0]} alt={data.data.name} fill={true} className='rounded-lg object-cover' />
                </div>
                {data.data.urls.length > 1 && <div className='w-[400px] overflow-auto h-full flex-col flex gap-4 max-h-[700px]'>
                    {data.data.urls.map((url) => <div key={url} className="h-[220px] relative w-full rounded-lg">
                        <Image src={url} alt={data.data.name} fill={true} className="object-cover rounded-lg" />
                    </div>)}
                </div>}
            </div>
            <div className='h-8 justify-end flex gap-6'>
                {
                    data.data.status === "verified" && <BsAward className="!text-4xl text-secondary" />
                }
                <FaExternalLinkAlt className='text-white/80 !text-lg' />
                <FaXTwitter className='text-white/80 !text-lg' />
                <FaDiscord className='text-white/80 !text-lg' />
            </div>
        </section>
        <section className="flex flex-col md:flex-row gap-6 h-full relative">
            <div className="md:flex-1 flex flex-col gap-6 h-full relative rounded-xl gradient p-8 text-[16px]">
                <div className="flex justify-between">
                    <div className="text-white">Contract Address</div>
                    <div className='text-white'>0x49cf...a28b</div>
                </div>
                <div className="flex justify-between">
                    <div className="text-white">Token ID</div>
                    <div className='text-white'>{data.data.tokenId}</div>
                </div><div className="flex justify-between">
                    <div className="text-white">Token Standard</div>
                    <div className='text-white'>{data.data.tokenStandard}</div>
                </div><div className="flex justify-between">
                    <div className="text-white">Blockchain</div>
                    <div className='text-white'>{data.data.blockchain}</div>
                </div><div className="flex justify-between">
                    <div className="text-white">Floor price</div>
                    <div className='text-white'>{data.data.floorPrice}</div>
                </div>
            </div>
            <Invest asset={data.data} />
        </section>
        <section>
            <ul className="flex gap-10">
                {navs.map(({ text, href }) => <li key={text}>
                    <Link href={href} className={(isLinkActive(href, path) ? "ring-1 ring-primary/40 rounded-md" : "") + ' p-1.5 text-[24px] text-white'}>{text}</Link>
                </li>)}
            </ul>
            <section className="mt-6">
                {children}
            </section>
        </section>
    </section> : <Error />)
}

export default Asset

const Invest = ({ asset }: { asset: TAsset }) => {
    const { user } = useContext(UserContext)
    const [quantity, setQuantity] = useState(1);
    const increaseQuantity = () => {
        if (quantity < asset.totalSupply) {
            setQuantity(quantity + 1);
        } else {
            toast.warn("Max supply reached!");
        }
    }
    const decreaseQuantity = () => {
        if (quantity > 1) {
            setQuantity(quantity - 1);
        } else {
            toast.warn("Min supply reached!");
        }
    }

    const [showBuyModal, setShowBuyModal] = useState(false);
    const [showSellModal, setShowSellModal] = useState(false);

    return <div className="flex-[2] ring-1 rounded-xl gradient ring-primary/40 flex flex-col items-center justify-center px-8 py-12 gap-8">
        {
            showBuyModal && <Modal onTapOutside={() => { setShowBuyModal(false) }}>
                <OrderModal asset={asset} setDefQuantity={setQuantity} mode="buy" defQuantity={quantity} onClose={() => setShowBuyModal(false)} />
            </Modal>
        }
        {
            showSellModal && <Modal onTapOutside={() => { setShowSellModal(false) }}>
                <OrderModal asset={asset} setDefQuantity={setQuantity} mode="sell" defQuantity={quantity} onClose={() => setShowSellModal(false)} />
            </Modal>
        }
        <progress className="w-full bg-secondary rounded-full" />
        <div className="flex w-full gap-24 items-center justify-evenly">
            <div className='flex flex-col flex-1 gap-4'>
                <div className='flex gap-2 justify-between items-center'>
                    <div className="text-white/40 text-sm">Expected Return</div>
                    <div className="text-white text-xl font-semibold">10.89%</div>
                </div>
                <div className='flex gap-2 justify-between items-center'>
                    <div className="text-white/40 text-sm">Investment Term</div>
                    <div className="text-white text-xl font-semibold">10 years</div>
                </div>
                <div className='flex gap-2 justify-between items-center'>
                    <div className="text-white/40 text-sm">Distribution Frequency</div>
                    <div className="text-white text-xl font-semibold">Yearly</div>
                </div>
            </div>
            <div className='flex-1 flex gap-4 items-center'>
                <span className='text-[24px] text-white'>Quantity</span>
                <div className="flex items-center gap-4 px-4 py-2 ring-1 ring-border rounded-md ">
                    <IoMdRemove className="text-white cursor-pointer" onClick={decreaseQuantity} /><span className="text-white select-none">{quantity}</span><IoMdAdd className="text-white cursor-pointer" onClick={increaseQuantity} />
                </div>
            </div>
        </div>
        <div className="w-full flex gap-8 justify-center">
            {
                (user.id === asset.creatorId && !asset.isListedByUser) && <button onClick={() => setShowSellModal(true)} className='flex-1 py-4 rounded-lg text-secondary ring-1 ring-primary/40'>Sell</button>
            }
            {
                user.id !== asset.creatorId && <button onClick={() => setShowBuyModal(true)} className="py-4 rounded-lg bg-secondary text-white text-sm w-full">Buy</button>
            }
        </div>
    </div>
}


const OrderModal = ({ asset, mode, defQuantity, setDefQuantity, onClose }: { asset: TAsset, mode: "buy" | "sell", defQuantity: number, setDefQuantity: React.Dispatch<React.SetStateAction<number>>, onClose: () => void }) => {
    const queryClient = new QueryClient();
    const [kind, setKind] = useState("limit")
    const [quantity, setQuantity] = useState(defQuantity);
    const [price, setPrice] = useState("");
    const increaseQuantity = () => {
        if (quantity < asset.totalSupply) {
            setQuantity(quantity + 1);
            setDefQuantity(quantity + 1);
        } else {
            toast.warn("Max supply reached!");
        }
    }
    const decreaseQuantity = () => {
        if (quantity > 1) {
            setQuantity(quantity - 1);
            setDefQuantity(quantity - 1);
        } else {
            toast.warn("Min supply reached!");
        }
    }
    const { mutateAsync, isPending } = OrderAPI.create({
        onSuccess: () => {
            setQuantity(1);
            onClose();
            queryClient.invalidateQueries({ queryKey: ["get order book"] });
        }
    })
    const handleClickInvest = async (type: "sell" | "buy") => {
        if (kind === "limit" && !price) {
            toast.error("Price is requried for limit order!");
            return;
        }
        await toast.promise(mutateAsync(kind === "limit" ? { assetId: asset.id, kind: kind, quantity, type: type, price: parseFloat(price) } : { assetId: asset.id, kind: kind, quantity, type: type }), {
            error: "Something went wrong, please try again!",
            pending: "Investing..., Please wait!",
            success: "Investment successful!"
        })
    }
    return <div className='w-[600px] max-h-[calc(100vh-100px)] overflow-auto z-20 rounded-lg bg-text-grey absolute roundd-md p-8 flex flex-col gap-4' onClick={(e) => e.stopPropagation()}>
        <ul className='ring-1 flex w-fit ring-white/40 '>
            {["Market", "Limit", "Stop"].map((val) => <li onClick={() => setKind(val.toLowerCase())} key={val} className={(val.toLowerCase() === kind ? '!bg-secondary !text-black' : '') + ' py-2 px-3 text-white text-[32px] cursor-pointer'}>{val}</li>)}
        </ul>
        <select className='ring-1 w-full px-2 py-3 rounded-lg text-sm text-white/80 bg-transparent ring-white/40'>
            <option>Good- Till Cancelled</option>
        </select>
        <div className='flex flex-col gap-2'>
            <div className='text-[18px] text-white/60'><span>Balance</span>: <span>0 USDT</span></div>
            <div className='text-[18px] text-white/60'><span>Best ask</span>: <span>0 USDT</span></div>
        </div>
        <div className='flex gap-4 items-center'>
            <span className='text-[24px] text-white'>Quantity</span>
            <div className="flex items-center gap-4 px-4 py-2 ring-1 ring-white/40 rounded-md ">
                <IoMdRemove className="text-white cursor-pointer" onClick={decreaseQuantity} /><span className="text-white select-none">{quantity}</span><IoMdAdd className="text-white cursor-pointer" onClick={increaseQuantity} />
            </div>
        </div>
        {
            kind === "limit" && <div className='flex w-full gap-8 items-center'>
                <span className='text-[24px] text-white'>Price</span>
                <input value={price} onChange={(e) => setPrice(e.target.value)} type='number' className="flex-1  px-4 py-3 bg-transparent ring-1 ring-white/40 rounded-md  text-white" />
            </div>
        }
        <div className='flex w-full gap-8 items-center'>
            <span className='text-[24px] text-white'>Total</span>
            <input disabled type='number' className="flex-1  px-4 py-3 bg-transparent ring-1 ring-white/40 rounded-md  text-white" />
        </div>
        <div className='flex flex-col gap-2 mt-8'>
            <div className='flex justify-between w-full items-center'>
                <span className='text-[18px] text-white/40'>Creator Royalties fee 0%</span>
                <span className='text-[18px] text-white/60'>0 USDT</span>
            </div>
            <div className='flex justify-between w-full items-center'>
                <span className='text-[18px] text-white/40'>Service fee 2.5%</span>
                <span className='text-[18px] text-white/60'>0.27 USDT</span>
            </div>
            <div className='flex justify-between w-full items-center'>
                <span className='text-[18px] text-white/40'>Total</span>
                <span className='text-[18px] text-white/60'>11.16 USDT</span>
            </div>
        </div>
        <button className='w-full mt-8 py-3 rounded-md bg-secondary text-[18px] capitalize' onClick={() => handleClickInvest(mode)} disabled={isPending}>{mode} Limit</button>
    </div>
}