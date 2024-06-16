export default function UserProfilePage(props: any) {
	return (
		<div class="bg-secondary p-4">
			<h1>Hello, {props.params.username}!</h1>
		</div>
	)
}
