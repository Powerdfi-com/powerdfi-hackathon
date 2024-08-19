"use client"
import Error from '@/app/components/error';
import Loading from '@/app/components/loading';
import AssetAPI from '@/app/utils/apis/asset'
import { dateFromISO, months } from '@/app/utils/func';
import React from 'react'

const AssetOrderBook = ({ params }: { params: { assetId: string, } }) => {
    const { isPending, data } = AssetAPI.getOrderBook({ id: params.assetId, status: "filled" });
    const { isPending: isPendingSell, data: sellOrders } = AssetAPI.getOrderBook({ id: params.assetId, status: "open", type: "sell" });
    const { isPending: isPendingBuy, data: buyOrders } = AssetAPI.getOrderBook({ id: params.assetId, status: "open", type: "buy" });
    return <section>
        <h4 className="text-secondary font-semibold text-[30px]">Open Orders</h4>
        <section className='w-full flex gap-16 mt-4 items-start'>
            <table className="text-white/80 table-fixed flex-1">
                <thead className="border-b border-primary/40 text-left  text-[24px] ">
                    <tr>
                        <th>Price (USDT)</th>
                        <th>Buy Order</th>
                    </tr>
                </thead>
                {
                    isPendingBuy ? <Loading /> : (buyOrders ? <tbody>
                        {
                            buyOrders.data.orders.map((order) => {
                                return <tr key={order.id} className="text-[20px] h-14">
                                    <td>{order.price}</td>
                                    <td>{order.quantity} Tokens</td>
                                </tr>
                            })
                        }
                    </tbody> : <Error />)
                }
            </table>
            <table className="text-white/80 table-fixed flex-1">
                <thead className="border-b border-primary/40 text-left  text-[24px] ">
                    <tr>
                        <th>Price (USDT)</th>
                        <th>Sell Order</th>
                    </tr>
                </thead>
                {
                    isPendingSell ? <Loading /> : (sellOrders ? <tbody>
                        {
                            sellOrders.data.orders.map((order) => {
                                return <tr key={order.id} className="text-[20px] h-14">
                                    <td>{order.price}</td>
                                    <td>{order.quantity} Tokens</td>
                                </tr>
                            })
                        }
                    </tbody> : <Error />)
                }
            </table>
        </section>
        <h4 className="text-secondary font-semibold text-[30px]  mt-16">Recently Filled Orders</h4>
        <table className="text-white/80 table-fixed w-full mt-4">
            <thead className="border-b border-primary/40 text-left  text-[24px] ">
                <tr>
                    <th>Filled</th>
                    <th>Quantity</th>
                    <th>Transaction Price</th>
                    <th>Type</th>
                </tr>
            </thead>
            {
                isPending ? <Loading /> : (data ? <tbody>
                    {
                        data.data.orders.map((order) => {
                            const date = dateFromISO(order.createdAt);
                            const thisYear = new Date().getFullYear();
                            const [day, month, year] = [date.getDay(), months[date.getMonth()], date.getFullYear()];
                            return <tr key={order.id} className="text-[20px] h-14">
                                <td>{`${month} ${day} ${thisYear === year && thisYear}`}</td>
                                <td>{order.quantity}</td>
                                <td>{order.price}</td>
                                <td>{order.type}</td>
                            </tr>
                        })
                    }
                </tbody> : <Error />)
            }
        </table>
    </section>
}

export default AssetOrderBook