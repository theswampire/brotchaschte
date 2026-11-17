<script lang="ts">
  import type { Order } from "../api/orders";

  interface Props {
    order: Order;
  }
  const { order: o }: Props = $props();
</script>

<td
  class="block max-w-20 overflow-hidden text-ellipsis px-2 py-1 border-b border-gray-400"
  >{o.id}</td
>
<td class="px-1 border-l border-b border-gray-400"
  >{o.firstname} {o.lastname}</td
>
<td class="px-1 border-l border-b border-gray-400">{o.email}</td>
<td class="px-1 border-l border-b border-gray-400">{o.phonenumber}</td>
<td class="px-1 border-l border-b border-gray-400">{o.address}</td>
<td class="px-1 border-l border-b border-gray-400">{o.postalcode}</td>
<td class="px-1 border-l border-b border-gray-400">{o.city}</td>
<td class="px-1 border-l border-b text-center border-gray-400"
  >{o.amountsizeS}</td
>
<td class="px-1 border-l border-b text-center border-gray-400"
  >{o.amountsizeM}</td
>
<td class="px-1 border-l border-b text-center border-gray-400"
  >{o.amountsizeL}</td
>
<td class="px-1 border-l border-b border-gray-400">{o.deliverymethod}</td>
<td class="px-1 border-l border-b border-gray-400">{o.remarks}</td>
